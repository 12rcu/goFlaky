package main

import (
	"goFlaky/adapters/junit"
	"goFlaky/adapters/persistence"
	"goFlaky/adapters/terminalui"
	"goFlaky/adapters/util"
	"goFlaky/core"
	"goFlaky/core/execution"
	"goFlaky/core/framework"
	"goFlaky/core/progress"
	"goFlaky/core/testmodify"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	config, err := core.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	for _, p := range config.Projects {
		var frameworkConfig framework.Config
		switch p.Framework {
		case "junit":
			frameworkConfig = junit.CreateNew()
		default:
			log.Fatalf("Unsupported framework: %s", p.Framework)
		}

		testFiles := util.SearchFileWithContent(
			p.ProjectDir+"/"+p.TestDir,
			func(path string, fileName string, fileContent string) bool {
				return strings.Contains(fileContent, frameworkConfig.Language())
			},
		)

		runOrders := [][]int{
			{0, 1, 2},
			{2, 1, 0},
		}
		testContent, err := os.ReadFile(testFiles[0])
		if err != nil {
			panic(err)
		}

		modifiedFiles := testmodify.ModifyTestFiles(runOrders, string(testContent), frameworkConfig)
		log.Println("modifiedFiles" + strconv.Itoa(len(modifiedFiles)))
		err = os.WriteFile(testFiles[0], []byte(modifiedFiles[0]), 0644)
		if err != nil {
			panic(err)
		}

		results := junit.ResultCollection(p.ProjectDir + "/" + p.TestResultDir)
		for _, result := range results {
			log.Println("Suite: " + result.TestSuite + " Test: " + result.TestName + " Outcome: " + result.TestOutcome)
		}
	}

	prgs := core.CreateProgressSlice(config)
	db, err := persistence.CreateSQLiteConnection()

	if err != nil {
		panic(err)
	}

	progressChannel := make(chan []progress.ProjectProgress)
	logChannel := make(chan string)
	go test(progressChannel, logChannel)

	dj := core.DependencyInjection{
		Config:             config,
		Progress:           prgs,
		ProgressChannel:    progressChannel,
		TerminalLogChannel: logChannel,
		Db:                 db,
	}

	go terminalui.TerminalUi(config, prgs, progressChannel, logChannel)
	execution.Execute(1, dj)
}

func test(progressChannel chan []progress.ProjectProgress, logChannel chan string) {
	for i := range 100 {
		time.Sleep(1 * time.Second)
		logChannel <- "Test " + strconv.Itoa(i)
	}
}
