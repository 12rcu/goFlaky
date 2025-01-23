package run

import (
	"goFlaky/adapters/mapper"
	"goFlaky/core"
	"goFlaky/core/execution"
	"log"
)

type IRunService interface {
	Execute() error
}

type Service struct {
	RunId int
	dj    core.DependencyInjection
}

func CreateService(runId int, dj core.DependencyInjection) *Service {
	return &Service{
		RunId: runId,
		dj:    dj,
	}
}

func (s *Service) Execute() {
	for _, p := range s.dj.Config.Projects {
		frameworkConfig, err := mapper.CreateFrameworkConfig(p.Framework)
		if err != nil {
			log.Printf("Error creating framework config %v", err)
		}

		workOrders := make(chan execution.WorkInfo)

		//pre runs
		preRunExecution := execution.PreRunExecution{
			RunId:      s.RunId,
			Project:    p,
			Dj:         s.dj,
			WorkOrders: workOrders,
		}
		odExecution := execution.OdRunExecution{
			RunId:         s.RunId,
			Project:       p,
			Dj:            s.dj,
			FrameworkConf: frameworkConfig,
			WorkOrders:    workOrders,
		}

		preRunExecution.ExecutePreRuns()
		odExecution.ExecuteOdRuns()
	}
}
