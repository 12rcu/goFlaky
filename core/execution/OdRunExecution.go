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

func (s OdRunExecution) ExecuteOdRuns() {
	testOrderStrategy, err := mapper.CreateStrategy(s.Project.Strategy)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, maxTestsPerFile := util.SearchFilesAndCount(
		s.Project.ProjectDir+"/"+s.Project.TestDir,
		s.FrameworkConf.TestAnnotation(),
	)

	orders := testOrderStrategy.GenerateOrder(maxTestsPerFile)

	for _, order := range orders {
		s.WorkOrders <- WorkInfo{
			RunId:     s.RunId,
			ProjectId: s.Project.Identifier,
			ModifyTestFile: func(projectDir string) {
				testFiles := util.SearchFileWithContent(
					projectDir,
					func(path string, fileName string, fileContent string) bool {
						return s.FrameworkConf.TestAnnotation().MatchString(fileContent)
					},
				)
				for _, file := range testFiles {
					content, err := os.ReadFile(file)
					if err != nil {
						log.Println(err.Error())
						continue
					}
					ov := testmodify.ModifyTestFile(order, string(content), s.FrameworkConf)
					err = os.WriteFile(file, []byte(ov), 0770)
					if err != nil {
						log.Println(err.Error())
					}
				}
			},
			GetTestOrder: func(testSuite string, testName string) []int { return order },
			Progress:     s.Dj.Progress,
		}
	}
}
