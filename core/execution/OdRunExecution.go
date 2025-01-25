package execution

import (
	"goFlaky/adapters/mapper"
	"goFlaky/adapters/util"
	"goFlaky/core"
	"goFlaky/core/framework"
	"goFlaky/core/testmodify"
	"log"
	"os"
)

type OdRunExecution struct {
	RunId         int
	Project       core.ConfigProject
	Dj            core.DependencyInjection
	FrameworkConf framework.Config
	WorkOrders    chan WorkInfo
}

func (exec OdRunExecution) runs() int {
	_, maxTestsPerFile := util.SearchFilesAndCount(
		exec.Project.ProjectDir+"/"+exec.Project.TestDir,
		exec.FrameworkConf.TestAnnotation(),
		func(log string) {
			exec.Dj.TerminalLogChannel <- "[ERROR] while searching for test orders " + log
		},
	)
	return maxTestsPerFile
}

func (exec OdRunExecution) ExecuteOdRuns() {
	testOrderStrategy, err := mapper.CreateStrategy(exec.Project.Strategy)
	if err != nil {
		log.Fatal(err.Error())
	}

	orders := testOrderStrategy.GenerateOrder(exec.runs())

	for _, order := range orders {
		exec.WorkOrders <- WorkInfo{
			RunId:     exec.RunId,
			ProjectId: exec.Project.Identifier,
			ModifyTestFile: func(projectDir string) {
				testFiles := util.SearchFileWithContent(
					projectDir,
					func(path string, fileName string, fileContent string) bool {
						return exec.FrameworkConf.TestAnnotation().MatchString(fileContent)
					},
					func(log string) {
						exec.Dj.TerminalLogChannel <- "[ERROR] while executing OD tests " + log
					},
				)
				for _, file := range testFiles {
					content, err := os.ReadFile(file)
					if err != nil {
						exec.Dj.TerminalLogChannel <- "[ERROR] OD failed to read test file " + err.Error()
						continue
					}
					ov := testmodify.ModifyTestFile(order, string(content), exec.FrameworkConf)
					err = os.WriteFile(file, []byte(ov), 0770)
					if err != nil {
						exec.Dj.TerminalLogChannel <- "[ERROR] OD failed to modify test file " + err.Error()
					}
				}
			},
			GetTestOrder: func(testSuite string, testName string) []int { return order },
			Progress:     exec.Dj.Progress,
		}
	}
}
