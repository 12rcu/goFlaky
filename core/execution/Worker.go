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
)

type WorkInfo struct {
	RunId          int
	ProjectId      string
	ModifyTestFile func(projectDir string)
	GetTestOrder   func(testSuite string, testName string) []int
	Progress       []progress.ProjectProgress
}

func (worker WorkInfo) WorkerExec(workerId int, project core.ConfigProject, dj core.DependencyInjection) {
	cleanAndCopyProjectToWorkingDir(workerId, project, dj)
	workDir, err := getWorkerDir(workerId, dj)
	if err != nil {
		log.Fatal(err)
		return
	}
	worker.ModifyTestFile(workDir)
	var testExecDir string
	if project.TestExecutionDir == "" {
		testExecDir = workDir
	} else {
		testExecDir = workDir + "/" + project.TestExecutionDir
	}
	OsCommand(project.TestExecutionCommand, testExecDir, dj.FileLogChannel, dj.TerminalLogChannel)
	worker.collectResultsFramework(project, dj)

	err = progress.ProgressProject(project.Identifier, worker.Progress)
	if err != nil {
		log.Println(err)
	}
	dj.ProgressChannel <- worker.Progress
}

func (worker WorkInfo) collectResultsFramework(project core.ConfigProject, dj core.DependencyInjection) {
	var path string
	if project.TestResultDir == "" {
		path = project.ProjectDir
	} else {
		path = project.ProjectDir + "/" + project.TestResultDir
	}
	switch project.Framework {
	case "jUnit":
		worker.collectResults(junit.ResultCollection(path), dj)
	default:
		log.Printf("Unsupported framework: %s", project.Framework)
	}
}

func (worker WorkInfo) collectResults(results []framework.TestResult, dj core.DependencyInjection) {
	for _, result := range results {
		var testOrder []string
		for _, o := range worker.GetTestOrder(result.TestSuite, result.TestName) {
			testOrder = append(testOrder, strconv.Itoa(o))
		}

		progressIndex, err := progress.GetProgressIndex(worker.ProjectId, worker.Progress)
		if err != nil {
			log.Println(err)
			return
		}
		err = persistence.CreateTestResult(
			dj.Db,
			worker.RunId,
			worker.Progress[progressIndex].Status,
			worker.ProjectId, result,
			strings.Join(testOrder, ","),
		)
		if err != nil {
			log.Println(err)
		}
	}
}

func cleanAndCopyProjectToWorkingDir(workerId int, project core.ConfigProject, dj core.DependencyInjection) {
	workerPath, err := getWorkerDir(workerId, dj)
	if err != nil {
		log.Default().Printf("Failed to get worker dir for %v: %v", workerId, err)
		return
	}
	err = os.RemoveAll(workerPath)
	if err != nil {
		log.Default().Printf("Failed to remove %v: %v", workerPath, err)
		return
	}
	err = os.Mkdir(workerPath, 0770)
	if err != nil {
		log.Default().Printf("Failed to create %v: %v", workerPath, err)
		return
	}
	err = os.CopyFS(workerPath, os.DirFS(project.ProjectDir))
	if err != nil {
		log.Default().Printf("Failed to copy %v: %v", workerPath, err)
		return
	}
}

func getWorkerDir(workerId int, dj core.DependencyInjection) (string, error) {
	path := dj.Config.TmpDir + "/workers/worker-" + strconv.Itoa(workerId)
	err := os.MkdirAll(path, 0770)
	if err != nil {
		return "", err
	}
	return path, nil
}
