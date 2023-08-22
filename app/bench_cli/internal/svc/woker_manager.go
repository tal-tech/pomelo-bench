package svc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/panjf2000/ants/v2"
	"pomelo_bench/app/bench/benchclient"
	"pomelo_bench/app/bench_cli/internal/config"
	"sync"
	"time"
)

type Woker struct {
	Address string
	Client  benchclient.Bench
}

type WorkManager struct {
	cfg        config.Config
	wokers     []Woker
	baseUid    uint64 // 学生base
	baseRoomId uint64 // 房间id base
}

func NewWorkManager(cfg config.Config) *WorkManager {
	return &WorkManager{
		cfg:     cfg,
		wokers:  make([]Woker, 0),
		baseUid: 10000,
	}
}

func (m *WorkManager) Add(w Woker) {
	m.wokers = append(m.wokers, w)
}

func (m *WorkManager) Connect(roomNumber int, roomSize int, channel int, callback func(string, *benchclient.StartPlanResponse, error)) (uid string) {

	uid = uuid.NewString()

	var (
		oneRoomNumber        = uint64(roomNumber / len(m.wokers))
		firstRoomLeaveNumber = roomNumber % len(m.wokers) // 未整除剩余的
		unix                 = time.Now().Unix()          // 放点时间戳 防止重复
	)

	// 组建房间请求数据
	requests := m.analysisStartPlanRequest(roomSize, channel, oneRoomNumber, firstRoomLeaveNumber, unix)

	wg := sync.WaitGroup{}

	pool, _ := ants.NewPoolWithFunc(len(requests), func(index interface{}) {

		i := index.(int64)

		res, err := m.wokers[i].Client.StartPlan(context.Background(), requests[i])

		if callback != nil {

			callback(m.wokers[i].Address, res, err)
		}

		wg.Done()
	})

	defer pool.Release()

	for i := 0; i < len(requests); i++ {

		wg.Add(1)

		// Submit tasks one by one.
		_ = pool.Invoke(int64(i))
	}

	wg.Wait()

	return uid
}

func (m *WorkManager) Each(fu func(woker Woker)) {

	for i := 0; i < len(m.wokers); i++ {
		fu(m.wokers[i])
	}

}

// 异步动作
func (m *WorkManager) EachAsync(fu func(woker Woker)) {
	wg := sync.WaitGroup{}

	for i := 0; i < len(m.wokers); i++ {

		wg.Add(1)

		go func(index int) {

			fu(m.wokers[index])

			wg.Done()

		}(i)
	}

	wg.Wait()
}

// 组建链接请求
func (m *WorkManager) analysisStartPlanRequest(roomSize int, channel int, oneRoomNumber uint64, firstRoomLeaveNumber int, unix int64) (res []*benchclient.StartPlanRequest) {
	for i := 0; i < len(m.wokers); i++ {

		number := oneRoomNumber
		if i < firstRoomLeaveNumber { // 剩余的部分多分配一个任务
			number++
		}

		if number == 0 {
			continue
		}

		res = append(res, &benchclient.StartPlanRequest{
			Plan: &benchclient.Plan{
				BaseUid:    m.baseUid,
				RoomNumber: number,
				RoomIdPre:  fmt.Sprintf("bench_%d_%d", m.baseRoomId, unix),
				RoomSize:   uint64(roomSize),
				Address:    m.cfg.PomeloAddress,
				ChannelId:  uint64(channel),
				Timeout:    uint64(m.cfg.Timeout), // 默认20秒连接超时时间
			},
		})

		m.baseUid += number * uint64(roomSize)
		m.baseRoomId++
	}

	return res
}
