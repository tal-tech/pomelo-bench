package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo_bench/app/bench/internal/service/planmanager"
	"pomelo_bench/app/bench/internal/svc"
	"pomelo_bench/pb/bench"
)

type ClosePlanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClosePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClosePlanLogic {
	return &ClosePlanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClosePlan 关闭压测计划
func (l *ClosePlanLogic) ClosePlan(in *bench.ClosePlanRequest) (*bench.ClosePlanResponse, error) {

	var (
		plans []*planmanager.Plan
	)

	if in.Uuid != nil { // 说明单发

		p, err := l.svcCtx.PlanManager.GetPlan(in.GetUuid())
		if err != nil {
			return nil, err
		}

		plans = []*planmanager.Plan{p}

	} else {
		plans = l.svcCtx.PlanManager.GetAllPlan()
	}

	for _, p := range plans {

		err := l.svcCtx.PlanManager.ClosePlan(p)
		if err != nil {
			return nil, err
		}
	}

	return &bench.ClosePlanResponse{}, nil
}
