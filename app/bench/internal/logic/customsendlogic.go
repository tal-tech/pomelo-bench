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

	err := l.svcCtx.PlanManager.GroupDo(in.Uuid, func(plan *planmanager.Plan) error {
		return plan.PlanCustomSend(l.ctx, in.Pool, in.Number, in.Limit, in.Duration)
	})

	return &bench.CustomSendResponse{}, err
}
