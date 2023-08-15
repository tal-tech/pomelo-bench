package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo_bench/app/bench/internal/logic/transform"
	"pomelo_bench/app/bench/internal/svc"
	"pomelo_bench/pb/bench"
)

type ListPlanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPlanLogic {
	return &ListPlanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListPlan 查询压测计划
func (l *ListPlanLogic) ListPlan(in *bench.ListPlanRequest) (*bench.ListPlanResponse, error) {

	plans := l.svcCtx.PlanManager.ListPlan()

	p := make([]*bench.PlanMonitor, 0, len(plans))

	for i := 0; i < len(plans); i++ {

		p = append(p, &bench.PlanMonitor{
			Uuid:      plans[i].Uuid,
			Plan:      plans[i].Cfg,
			Status:    plans[i].Status,
			Connector: plans[i].Connector,
			Total:     transform.Statistics(plans[i].TotalStatistics),
			Stat:      transform.Stat(plans[i].Stat),
		})
	}

	return &bench.ListPlanResponse{Plans: p}, nil
}
