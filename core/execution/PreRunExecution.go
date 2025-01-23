package execution

import "goFlaky/core"

type PreRunExecution struct {
	RunId      int
	Project    core.ConfigProject
	Dj         core.DependencyInjection
	WorkOrders chan WorkInfo
}

func (s PreRunExecution) ExecutePreRuns() {
	for i := 0; i < int(s.Project.PreRuns); i++ {
		s.WorkOrders <- WorkInfo{
			RunId:          s.RunId,
			ProjectId:      s.Project.Identifier,
			ModifyTestFile: func(projectDir string) {},
			GetTestOrder:   func(testSuite string, testName string) []int { return []int{} },
			Progress:       s.Dj.Progress,
		}
	}
}
