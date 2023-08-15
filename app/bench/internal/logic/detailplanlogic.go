package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo_bench/app/bench/internal/logic/transform"
	"pomelo_bench/app/bench/internal/svc"
	"pomelo_bench/pb/bench"
)

type DetailPlanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailPlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailPlanLogic {
	return &DetailPlanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DetailPlan 查询压测计划详情
func (l *DetailPlanLogic) DetailPlan(in *bench.DetailPlanRequest) (*bench.DetailPlanResponse, error) {
	p, err := l.svcCtx.PlanManager.GetPlan(in.GetUuid())
	if err != nil {
		return nil, err
	}

	detail := p.PlanDetail(l.ctx)

	res := &bench.DetailPlanResponse{
		Uuid: in.GetUuid(),
		Detail: &bench.PlanDetail{
			Plan:       detail.Cfg,
			Connectors: transform.Connectors(detail.Connectors),
		},
		Status: detail.Status,
	}

	return res, nil
}
