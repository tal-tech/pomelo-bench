package logic

import (
	"context"
	"pomelo_bench/app/bench/internal/service/planmanager"

	"pomelo_bench/app/bench/internal/svc"
	"pomelo_bench/pb/bench"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearStatisticsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearStatisticsLogic {
	return &ClearStatisticsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClearStatistics 清理任务指标
func (l *ClearStatisticsLogic) ClearStatistics(in *bench.ClearStatisticsRequest) (*bench.ClearStatisticsResponse, error) {

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

		p.ClearMetrics(l.ctx)
	}

	return &bench.ClearStatisticsResponse{}, nil
}
