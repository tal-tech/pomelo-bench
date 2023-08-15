package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo_bench/app/bench/internal/logic/transform"
	"pomelo_bench/app/bench/internal/svc"
	"pomelo_bench/pb/bench"
)

type StartPlanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStartPlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartPlanLogic {
	return &StartPlanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// StartPlan 快速启动压测
func (l *StartPlanLogic) StartPlan(in *bench.StartPlanRequest) (*bench.StartPlanResponse, error) {
	// 创建任务
	uid := l.svcCtx.PlanManager.CreatePlan(in.Plan)

	plan, err := l.svcCtx.PlanManager.GetPlan(uid)
	if err != nil {
		return nil, err
	}

	l.Info("创建任务成功,准备通过网关链接获取chat地址")

	// 通过网关链接获取chat地址
	err = plan.PlanQueryGate(l.ctx)
	if err != nil {
		return nil, err
	}

	l.Info("通过网关链接获取chat地址成功,准备进入chat房间")

	// 进入chat房间
	err = plan.PlanConnectEntry(l.ctx)
	if err != nil {
		return nil, err
	}

	l.Info("进入chat房间成功")

	detail := plan.PlanDetail(l.ctx)

	res := &bench.StartPlanResponse{
		Uuid: uid,
		Detail: &bench.PlanDetail{
			Plan:       detail.Cfg,
			Connectors: transform.Connectors(detail.Connectors),
		},
		Status: detail.Status,
	}

	return res, nil
}
