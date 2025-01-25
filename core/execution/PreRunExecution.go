package execution

import "goFlaky/core"

type PreRunExecution struct {
	RunId      int
	Project    core.ConfigProject
	Dj         core.DependencyInjection
	WorkOrders chan WorkInfo
}

func (exec PreRunExecution) runs() int {
	return int(exec.Project.PreRuns)
}

func (exec PreRunExecution) ExecutePreRuns() {
	for i := 0; i < exec.runs(); i++ {
		exec.WorkOrders <- WorkInfo{
			RunId:          exec.RunId,
			ProjectId:      exec.Project.Identifier,
			ModifyTestFile: func(projectDir string) {},
			GetTestOrder:   func(testSuite string, testName string) []int { return []int{} },
			Progress:       exec.Dj.Progress,
		}
	}
}
