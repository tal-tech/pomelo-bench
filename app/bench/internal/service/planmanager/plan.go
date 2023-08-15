package planmanager

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/timex"
	"pomelo_bench/app/bench/internal/metrics"
	lcpomelo2 "pomelo_bench/app/bench/internal/service/lcpomelo"
	"pomelo_bench/pb/bench"
	"sync"
	"time"
)

type Plan struct {
	uid             string
	cfg             *bench.Plan
	status          bench.Status
	lcPomeloClients []*lcpomelo2.ClientConnector
	metrics         *metrics.Metrics
}

func NewPlan(uid string, cfg *bench.Plan) *Plan {

	p := &Plan{
		uid:             uid,
		cfg:             cfg,
		status:          bench.Status_Waiting,
		lcPomeloClients: make([]*lcpomelo2.ClientConnector, 0, cfg.RoomNumber*cfg.RoomSize),
		metrics:         metrics.NewMetrics("plan"),
	}

	for i := 0; i < int(cfg.RoomNumber); i++ {
		var roomId = fmt.Sprintf("%s_%d", cfg.RoomIdPre, i) // room id

		for j := 0; j < int(cfg.RoomSize); j++ {

			var uid = int(cfg.BaseUid) + i*int(cfg.RoomSize) + j // 学员id

			c := lcpomelo2.NewClientConnector(cfg.Address, uid, int(cfg.ChannelId), roomId)

			p.lcPomeloClients = append(p.lcPomeloClients, c)
		}
	}

	return p
}

func (p *Plan) PlanQueryGate(ctx context.Context) error {

	err := p.asyncDo(func(connector *lcpomelo2.ClientConnector) error {

		err := connector.RunGateConnectorAndWaitConnect(ctx)
		if err != nil {

			p.status = bench.Status_Failed

			return err
		}

		err = connector.SyncGateRequest(ctx)
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

func (p *Plan) PlanConnectEntry(ctx context.Context) error {

	err := p.asyncDo(func(connector *lcpomelo2.ClientConnector) error {

		err := connector.RunChatConnectorAndWaitConnect(ctx)
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

	return err
}

func (p *Plan) PlanSendChat(ctx context.Context, message string, number uint64, limit uint64, duration uint64) (err error) {
	if number == 0 {
		return nil
	}

	err = p.asyncDo2(int(limit), func(index int, length int, connector *lcpomelo2.ClientConnector) error {

		for j := 0; j < int(number); j++ {

			startTime := timex.Now()
			err := connector.AsyncChatSendMessage(ctx, message, func(data []byte) {
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

	err = p.asyncDo2(int(limit), func(index int, length int, connector *lcpomelo2.ClientConnector) error {

		for i := 0; i < int(number); i++ {

			for j := index; j < len(pool.Data); j += length { // 只取自己的那条子集

				startTime := timex.Now()
				err := connector.AsyncCustomSend(ctx, pool.Router, pool.Data[j], func(data []byte) {
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
		}

		return nil
	})

	return err
}

func (p *Plan) PlanDetail(ctx context.Context) PlanDetail {

	res := PlanDetail{
		Cfg:        p.cfg,
		Connectors: make([]lcpomelo2.ClientDetail, 0, len(p.lcPomeloClients)),
		Status:     p.status,
		Stat:       p.metrics.Execute(),
	}

	for i := 0; i < len(p.lcPomeloClients); i++ {
		detail := p.lcPomeloClients[i].Detail(ctx)

		res.Connectors = append(res.Connectors, detail)
	}

	return res
}

func (p *Plan) Close(ctx context.Context) error {

	err := p.asyncDo(func(connector *lcpomelo2.ClientConnector) error {
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

// 异步动作
func (p *Plan) asyncDo(do func(*lcpomelo2.ClientConnector) error) (err error) {

	wg := sync.WaitGroup{}

	for i := 0; i < len(p.lcPomeloClients); i++ {

		wg.Add(1)

		go func(index int) {

			oErr := do(p.lcPomeloClients[index])
			if oErr != nil {
				err = oErr
			}

			wg.Done()

		}(i)
	}

	wg.Wait()

	return err
}

// 异步动作2
func (p *Plan) asyncDo2(limit int, do func(index int, length int, connector *lcpomelo2.ClientConnector) error) (err error) {
	l := len(p.lcPomeloClients)
	if limit != 0 && l > limit {
		l = limit
	}

	wg := sync.WaitGroup{}

	for i := 0; i < l; i++ {

		wg.Add(1)

		go func(index int) {

			oErr := do(index, l, p.lcPomeloClients[index])
			if oErr != nil {
				err = oErr
			}

			wg.Done()

		}(i)
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
