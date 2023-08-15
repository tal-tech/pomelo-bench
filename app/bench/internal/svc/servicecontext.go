package svc

import (
	"pomelo_bench/app/bench/internal/config"
	"pomelo_bench/app/bench/internal/service/planmanager"
)

type ServiceContext struct {
	Config      config.Config
	PlanManager *planmanager.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {

	planManager := planmanager.NewManager()

	return &ServiceContext{
		Config: c,

		PlanManager: planManager,
	}
}
