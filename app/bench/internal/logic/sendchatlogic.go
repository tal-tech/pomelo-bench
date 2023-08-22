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

	err := l.svcCtx.PlanManager.GroupDo(in.Uuid, func(plan *planmanager.Plan) error {
		return plan.PlanSendChat(l.ctx, in.Message, in.Number, in.Limit, in.Duration)
	})

	return &bench.SendChatResponse{}, err
}
