package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo_bench/app/bench/internal/service/planmanager"
	"pomelo_bench/app/bench/internal/svc"
	"pomelo_bench/pb/bench"
)

type SendChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendChatLogic {
	return &SendChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendChat 快速发送消息
func (l *SendChatLogic) SendChat(in *bench.SendChatRequest) (*bench.SendChatResponse, error) {

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
		err := p.PlanSendChat(l.ctx, in.Message, in.Number, in.Limit, in.Duration)
		if err != nil {
			return nil, err
		}
	}

	return &bench.SendChatResponse{}, nil
}
