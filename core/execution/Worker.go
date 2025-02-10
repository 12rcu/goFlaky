package execution

import (
	"goFlaky/adapters/junit"
	"goFlaky/core"
	"goFlaky/core/framework"
	"goFlaky/core/progress"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type WorkInfo struct {
	RunId          int
	ProjectId      string
	ModifyTestFile func(projectDir string)
	GetTestOrder   func(testSuite string, testName string) []int
	Progress       []progress.ProjectProgress
	RunType        string
}

func Worker(id int, project core.ConfigProject, dj core.DependencyInjection, jobs <-chan WorkInfo, waitGroup *sync.WaitGroup) {
	for j := range jobs {
		j.workerExec(id, project, dj)
	}
	waitGroup.Done()
}

func (order WorkInfo) workerExec(workerId int, project core.ConfigProject, dj core.DependencyInjection) {
	//copy project to tmp dir
	cleanAndCopyProjectToWorkingDir(workerId, project, dj)
	workDir, err := getWorkerDir(workerId, dj)
	if err != nil {
		log.Fatal(err)
		return
	}
	//modify test files if needed
	order.ModifyTestFile(workDir)
	var testExecDir string
	if project.TestExecutionDir == "" {
		testExecDir = workDir
	} else {
		testExecDir = workDir + "/" + project.TestExecutionDir
	}

	//run tests
	OsCommand(project.TestExecutionCommand, testExecDir, dj.FileLogChannel, dj.TerminalLogChannel)
	//collect results
	order.collectResultsFramework(project, dj, workDir)

	//update progress for ui feedback
	err = progress.ProgressProject(project.Identifier, order.Progress)
	if err != nil {
		dj.TerminalLogChannel <- "[ERROR] worker, " + err.Error()
	}
	dj.ProgressChannel <- order.Progress
}

func (order WorkInfo) collectResultsFramework(project core.ConfigProject, dj core.DependencyInjection, workDir string) {
	var path string
	if project.TestResultDir == "" {
		path = workDir
	} else {
		path = workDir + "/" + project.TestResultDir
	}
	switch strings.ToLower(project.Framework) {
	case "junit":
		order.collectResults(junit.ResultCollection(path, dj), dj)
		break
	default:
		dj.TerminalLogChannel <- "[ERROR] Framework not supported: " + project.Framework
	}
}

func (order WorkInfo) collectResults(results []framework.TestResult, dj core.DependencyInjection) {
	for _, result := range results {
		var testOrder []string
		for _, o := range order.GetTestOrder(result.TestSuite, result.TestName) {
			testOrder = append(testOrder, strconv.Itoa(o))
		}
		if order.RunType == "OD_RUN" && len(testOrder) <= 0 {
			return //disabled test
		}

		dj.TerminalLogChannel <- "Added test result for [" + result.TestName + "]"
		err := dj.Db.CreateTestResult(
			order.RunId,
			order.RunType,
			order.ProjectId, result,
			strings.Join(testOrder, ","),
		)
		if err != nil {
			dj.TerminalLogChannel <- "[ERROR] " + err.Error()
		}
	}
}

func cleanAndCopyProjectToWorkingDir(workerId int, project core.ConfigProject, dj core.DependencyInjection) {
	workerPath, err := getWorkerDir(workerId, dj)
	if err != nil {
		dj.TerminalLogChannel <- "[ERROR] Failed to get worker dir " + err.Error()
		return
	}
	err = os.RemoveAll(workerPath)
	if err != nil {
		dj.TerminalLogChannel <- "[ERROR] Failed to clear work dir " + err.Error()
		return
	}
	err = os.Mkdir(workerPath, 0777)
	if err != nil {
		dj.TerminalLogChannel <- "[ERROR] Failed to create work dir " + err.Error()
		return
	}
	err = os.CopyFS(workerPath, os.DirFS(project.ProjectDir))
	if err != nil {
		dj.TerminalLogChannel <- "[ERROR] Failed to copy to work dir " + err.Error()
		return
	}
}

func getWorkerDir(workerId int, dj core.DependencyInjection) (string, error) {
	path := dj.Config.BaseDir + "/" + dj.Config.TmpDir + "/workers/worker-" + strconv.Itoa(workerId)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return "", err
	}
	return path, nil
}
