package planmanager

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"math/rand"
	"pomelo_bench/app/bench/internal/metrics"
	lcpomelo2 "pomelo_bench/app/bench/internal/service/lcpomelo"
	"pomelo_bench/pb/bench"
	"sync"
	"time"
)

const DefaultConcurrency = 100

type Plan struct {
	uid         string
	cfg         *bench.Plan
	status      bench.Status
	rooms       []room
	clientCount int
	metrics     *metrics.Metrics
}

type room struct {
	roomId  string
	clients []*lcpomelo2.ClientConnector
}

func NewPlan(uid string, cfg *bench.Plan) *Plan {

	p := &Plan{
		uid:         uid,
		cfg:         cfg,
		status:      bench.Status_Waiting,
		rooms:       make([]room, 0, cfg.RoomNumber),
		clientCount: int(cfg.RoomSize * cfg.RoomNumber),
		metrics:     metrics.NewMetrics("plan"),
	}

	for i := 0; i < int(cfg.RoomNumber); i++ {

		r := room{
			//roomId:  fmt.Sprintf("%s_%d", uuid.NewString(), i), // room id
			clients: make([]*lcpomelo2.ClientConnector, 0, cfg.RoomSize),
		}

		if len(cfg.RoomIds) == int(cfg.RoomNumber) {
			r.roomId = cfg.RoomIds[i]
		} else if cfg.RoomIdPre != nil {
			r.roomId = fmt.Sprintf("%s_%d", *cfg.RoomIdPre, i)
		} else {
			r.roomId = fmt.Sprintf("%s_%d", uuid.NewString(), i)
		}

		for j := 0; j < int(cfg.RoomSize); j++ {

			var uid = int(cfg.BaseUid) + i*int(cfg.RoomSize) + j // 学员id

			c := lcpomelo2.NewClientConnector(cfg.Address, uid, int(cfg.ChannelId), r.roomId)

			r.clients = append(r.clients, c)
		}

		p.rooms = append(p.rooms, r)
	}

	return p
}

func (p *Plan) PlanQueryGateAndEnter(ctx context.Context, timeout time.Duration) error {

	// 限制下并发量 20
	err := p.asyncDo(0, DefaultConcurrency, func(connector *lcpomelo2.ClientConnector) error {

		err := connector.RunGateConnectorAndWaitConnect(ctx, timeout)
		if err != nil {

			p.status = bench.Status_Failed

			return err
		}

		err = connector.SyncGateRequest(ctx)
		if err != nil {

			p.status = bench.Status_Failed

			return err
		}

		// 关闭gate失败 不认为失败
		err = connector.CloseGate(ctx)
		if err != nil {
			logx.Error("connector.CloseGate(ctx) failed,err:", err)
		}

		err = connector.RunChatConnectorAndWaitConnect(ctx, timeout)
		if err != nil {

			p.status = bench.Status_Failed

			return err
		}

		err = connector.SyncChatEnterConnectorRequest(ctx)
		if err != nil {

			p.status = bench.Status_Failed

			return err
		}

		return nil
	})

	if err != nil {
		return err
	} else {
		p.status = bench.Status_Doing
	}

	return nil
}

func (p *Plan) PlanSendChat(ctx context.Context, message string, number uint64, limit uint64, duration uint64) (err error) {
	if number == 0 {
		return nil
	}

	err = p.asyncDo(int(limit), 0, func(connector *lcpomelo2.ClientConnector) error {

		for j := 0; j < int(number); j++ {

			startTime := timex.Now()
			err := connector.AsyncChatSendMessage(ctx, message, func(_ []byte) {
				duration := timex.Since(startTime)
				p.metrics.Add(duration)
			})
			if err != nil {
				p.metrics.Drop()
				return err
			}

			if duration > 0 {
				time.Sleep(time.Duration(duration) * time.Millisecond)
			}
		}

		return nil
	})

	return err
}

func (p *Plan) PlanCustomSend(ctx context.Context, pool *bench.CustomMessagePool, number uint64, limit uint64, duration uint64) (err error) {
	if number == 0 {
		return nil
	}

	if len(pool.Data) == 0 {
		return nil
	}

	err = p.asyncDo(int(limit), 0, func(connector *lcpomelo2.ClientConnector) error {

		for i := 0; i < int(number); i++ {

			index := rand.Intn(len(pool.Data))

			data := pool.Data[index]

			startTime := timex.Now()
			err := connector.AsyncCustomSend(ctx, pool.Router, data, func(_ []byte) {
				duration := timex.Since(startTime)
				p.metrics.Add(duration)
			})
			if err != nil {
				p.metrics.Drop()
				return err
			}

			if duration > 0 {
				time.Sleep(time.Duration(duration) * time.Millisecond)
			}
		}

		return nil
	})

	return err
}

func (p *Plan) PlanDetail(ctx context.Context) PlanDetail {

	res := PlanDetail{
		Cfg:        p.cfg,
		Connectors: make([]lcpomelo2.ClientDetail, 0, p.clientCount),
		Status:     p.status,
		Stat:       p.metrics.Execute(),
	}

	for i := 0; i < len(p.rooms); i++ {

		for j := 0; j < len(p.rooms[i].clients); j++ {
			detail := p.rooms[i].clients[j].Detail(ctx)

			res.Connectors = append(res.Connectors, detail)
		}
	}

	return res
}

func (p *Plan) CloseGate(ctx context.Context) error {

	err := p.asyncDo(0, DefaultConcurrency, func(connector *lcpomelo2.ClientConnector) error {

		return connector.CloseGate(ctx)
	})

	return err
}

func (p *Plan) Close(ctx context.Context) error {

	err := p.asyncDo(0, DefaultConcurrency, func(connector *lcpomelo2.ClientConnector) error {
		err := connector.CloseGate(ctx)
		if err != nil {
			return err
		}

		err = connector.CloseChat(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (p *Plan) ClearMetrics(ctx context.Context) {
	p.metrics.Clear()

	for i := 0; i < len(p.rooms); i++ {

		for j := 0; j < len(p.rooms[i].clients); j++ {

			p.rooms[i].clients[j].ClearMetrics()

		}
	}

}

// 异步动作2 limit 限制房间发送人数 concurrency 并发量
func (p *Plan) asyncDo(roomLimit int, concurrency int, do func(connector *lcpomelo2.ClientConnector) error) (err error) {
	l := p.clientCount
	limitLen := roomLimit * int(p.cfg.RoomNumber)
	if limitLen != 0 && l > limitLen {
		l = limitLen
	}

	// 并发量最大等于 发送学员数
	size := l
	if concurrency != 0 && size > concurrency {
		size = concurrency
	}

	wg := sync.WaitGroup{}

	pool, _ := ants.NewPoolWithFunc(size, func(i interface{}) {
		ii := i.(int64)
		roomIndex := int(ii) % int(p.cfg.RoomNumber)
		clientIndex := int(ii) / int(p.cfg.RoomNumber)

		oErr := do(p.rooms[roomIndex].clients[clientIndex])
		if oErr != nil {
			err = oErr
		}

		wg.Done()
	})
	defer pool.Release()

	for i := 0; i < l; i++ {

		wg.Add(1)

		// Submit tasks one by one.
		_ = pool.Invoke(int64(i))
	}

	wg.Wait()

	return err
}

type PlanDetail struct {
	Cfg        *bench.Plan              // 计划配置
	Connectors []lcpomelo2.ClientDetail // 连接器详情
	Status     bench.Status             // 计划状态
	Stat       metrics.StatReport       // 999指标
}
