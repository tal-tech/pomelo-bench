package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo_bench/app/bench/internal/service/planmanager"
	"pomelo_bench/app/bench/internal/svc"
	"pomelo_bench/pb/bench"
)

type CustomSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCustomSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CustomSendLogic {
	return &CustomSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CustomSend 自定义消息发送
func (l *CustomSendLogic) CustomSend(in *bench.CustomSendRequest) (*bench.CustomSendResponse, error) {

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
		err := p.PlanCustomSend(l.ctx, in.Pool, in.Number, in.Limit, in.Duration)
		if err != nil {
			return nil, err
		}
	}

	return &bench.CustomSendResponse{}, nil
}
