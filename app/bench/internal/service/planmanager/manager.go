// Package planmanager 计划管理
package planmanager

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo_bench/app/bench/internal/metrics"
	"pomelo_bench/app/bench/internal/service/lcpomelo"
	"pomelo_bench/pb/bench"
	"sync"
)

type Manager struct {
	plans map[string]*Plan
	mu    sync.Mutex
}

func NewManager() *Manager {

	return &Manager{
		plans: make(map[string]*Plan, 0),
	}

}

func (m *Manager) CreatePlan(cfg *bench.Plan) string {
	uid := uuid.NewString()

	p := NewPlan(uid, cfg)

	m.plans[uid] = p

	return uid
}

func (m *Manager) ListPlan() (infos []PlanInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	infos = make([]PlanInfo, 0, len(m.plans))

	for uid, plan := range m.plans {

		detail := plan.PlanDetail(context.Background())
		var (
			// 总发送量
			totalStatistics lcpomelo.Statistics
			// Connector 客户端链接情况
			connector bench.ConnectorStatus
		)

		for i := 0; i < len(detail.Connectors); i++ {

			totalStatistics.Add(detail.Connectors[i].Statistics)

			if detail.Connectors[i].PomeloGate.ConnectorConnected {
				connector.GateConnector++
			}
			if detail.Connectors[i].PomeloChat.ConnectorConnected {
				connector.ChatConnector++
			}

		}

		infos = append(infos, PlanInfo{
			Uuid:            uid,
			Cfg:             plan.cfg,
			Status:          plan.status,
			Connector:       &connector,
			TotalStatistics: totalStatistics,
			Stat:            detail.Stat,
		})

	}

	return infos

}

func (m *Manager) GetPlan(uuid string) (*Plan, error) {

	p, ok := m.plans[uuid]
	if !ok {
		return nil, errors.New("invalid uuid")
	}

	return p, nil
}

func (m *Manager) GetAllPlan() []*Plan {

	res := make([]*Plan, 0, len(m.plans))
	for _, plan := range m.plans {
		res = append(res, plan)
	}

	return res
}

func (m *Manager) ClosePlan(p *Plan) error {

	err := p.Close(context.Background())
	if err != nil {

		logx.Error("QuickClosePlan failed ,err:", err)
		return err
	}

	delete(m.plans, p.uid)

	return nil
}

type PlanInfo struct {
	Uuid   string
	Cfg    *bench.Plan
	Status bench.Status

	// Connector 客户端链接情况
	Connector *bench.ConnectorStatus
	// 总发送量 指标统计
	TotalStatistics lcpomelo.Statistics
	// 指标统计
	Stat metrics.StatReport
}
