package execution

import (
	"goFlaky/adapters/junit"
	"goFlaky/adapters/persistence"
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
}

func Worker(id int, project core.ConfigProject, dj core.DependencyInjection, jobs <-chan WorkInfo, waitGroup *sync.WaitGroup) {
	for j := range jobs {
		j.workerExec(id, project, dj)
	}
	waitGroup.Done()
}

func (order WorkInfo) workerExec(workerId int, project core.ConfigProject, dj core.DependencyInjection) {
	cleanAndCopyProjectToWorkingDir(workerId, project, dj)
	workDir, err := getWorkerDir(workerId, dj)
	if err != nil {
		log.Fatal(err)
		return
	}
	order.ModifyTestFile(workDir)
	var testExecDir string
	if project.TestExecutionDir == "" {
		testExecDir = workDir
	} else {
		testExecDir = workDir + "/" + project.TestExecutionDir
	}
	OsCommand(project.TestExecutionCommand, testExecDir, dj.FileLogChannel, dj.TerminalLogChannel)
	order.collectResultsFramework(project, dj)

	err = progress.ProgressProject(project.Identifier, order.Progress)
	if err != nil {
		dj.TerminalLogChannel <- "[ERROR] worker, " + err.Error()
	}
	dj.ProgressChannel <- order.Progress
}

func (order WorkInfo) collectResultsFramework(project core.ConfigProject, dj core.DependencyInjection) {
	var path string
	if project.TestResultDir == "" {
		path = project.ProjectDir
	} else {
		path = project.ProjectDir + "/" + project.TestResultDir
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

		progressIndex, err := progress.GetProgressIndex(order.ProjectId, order.Progress)
		if err != nil {
			dj.TerminalLogChannel <- "[ERROR] " + err.Error()
			return
		}
		err = persistence.CreateTestResult(
			dj.Db,
			order.RunId,
			order.Progress[progressIndex].Status,
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
