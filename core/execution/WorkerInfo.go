package execution

import "goFlaky/core/progress"

type ExecInfo struct {
	RunId          int
	ProjectId      string
	ModifyTestFile func(projectDir string)
	GetTestOrder   func(testSuite string, testName string) []int
	Progress       []progress.ProjectProgress
}
