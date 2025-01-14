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

func WorkerExec(workerId int, project core.ConfigProject, dj core.DependencyInjection, info ExecInfo) {
	cleanAndCopyProjectToWorkingDir(workerId, project, dj)
	workDir, err := getWorkerDir(workerId, dj)
	if err != nil {
		log.Fatal(err)
		return
	}
	info.ModifyTestFile(workDir)
	var testExecDir string
	if project.TestExecutionDir == "" {
		testExecDir = workDir
	} else {
		testExecDir = workDir + "/" + project.TestExecutionDir
	}
	OsCommand(project.TestExecutionCommand, testExecDir, dj.FileLogChannel, dj.TerminalLogChannel)
	collectResultsFramework(project, dj, info)

	err = progress.ProgressProject(project.Identifier, info.Progress)
	if err != nil {
		log.Println(err)
	}
	dj.ProgressChannel <- info.Progress
}

func collectResultsFramework(project core.ConfigProject, dj core.DependencyInjection, info ExecInfo) {
	var path string
	if project.TestResultDir == "" {
		path = project.ProjectDir
	} else {
		path = project.ProjectDir + "/" + project.TestResultDir
	}
	switch project.Framework {
	case "jUnit":
		collectResults(junit.ResultCollection(path), dj, info)
	default:
		log.Printf("Unsupported framework: %s", project.Framework)
	}
}

func collectResults(results []framework.TestResult, dj core.DependencyInjection, info ExecInfo) {
	for _, result := range results {
		var testOrder []string
		for _, o := range info.GetTestOrder(result.TestSuite, result.TestName) {
			testOrder = append(testOrder, strconv.Itoa(o))
		}

		progressIndex, err := progress.GetProgressIndex(info.ProjectId, info.Progress)
		if err != nil {
			log.Println(err)
			return
		}
		err = persistence.CreateTestResult(
			dj.Db,
			info.RunId,
			info.Progress[progressIndex].Status,
			info.ProjectId, result,
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
