package main

import (
	"goFlaky/adapters/junit"
	"goFlaky/adapters/persistence"
	"goFlaky/adapters/terminalui"
	"goFlaky/core"
	"goFlaky/core/execution"
	"goFlaky/core/progress"
	"log"
	"strconv"
	"time"
)

func main() {
	config, err := core.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	for _, p := range config.Projects {
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
