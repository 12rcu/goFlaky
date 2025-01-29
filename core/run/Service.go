package run

import (
	"goFlaky/adapters/mapper"
	"goFlaky/core"
	"goFlaky/core/execution"
	"goFlaky/core/framework"
	"sync"
)

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

func (s *Service) Execute(waitGroup *sync.WaitGroup) {
	for _, p := range s.dj.Config.Projects {
		frameworkConfig, err := mapper.CreateFrameworkConfig(p.Framework)
		if err != nil {
			s.dj.TerminalLogChannel <- "[ERROR] while creating framework config " + err.Error()
		}

		workOrders := make(chan execution.WorkInfo)
		var workerWaitGroup sync.WaitGroup

		//schedule test runs with pre runs and od runs
		go s.scheduleWorkOrders(workOrders, frameworkConfig, p)

		//create n runners defined by the config
		for i := 0; i < int(s.dj.Config.Worker); i++ {
			workerWaitGroup.Add(1)
			//execute scheduled orders
			go execution.Worker(i, p, s.dj, workOrders, &workerWaitGroup)
		}

		workerWaitGroup.Wait()

		//results, err := s.dj.Db.GetProjectTestResults(s.RunId, p.Identifier)
		//if err != nil {
		//	return
		//}
		//execution.CreateClassification(results, s.dj, s.RunId, p.Identifier)
	}

	close(s.dj.ProgressChannel)
	close(s.dj.TerminalLogChannel)
	close(s.dj.FileLogChannel)

	waitGroup.Done()
}

func (s *Service) scheduleWorkOrders(workOrders chan execution.WorkInfo, frameworkConfig framework.Config, p core.ConfigProject) {
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
	close(workOrders)
}
