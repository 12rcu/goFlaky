package main

import (
	"goFlaky/adapters/persistence"
	"goFlaky/adapters/terminalui"
	"goFlaky/core"
	"goFlaky/core/progress"
	"goFlaky/core/run"
	"log"
	"strconv"
	"time"
)

func main() {
	config, err := core.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	prgs := core.CreateProgressSlice(config)
	db, err := persistence.CreateSQLiteConnection()

	if err != nil {
		log.Fatal(err)
	}

	progressChannel := make(chan []progress.ProjectProgress)
	logChannel := make(chan string)

	dj := core.DependencyInjection{
		Config:             config,
		Progress:           prgs,
		ProgressChannel:    progressChannel,
		TerminalLogChannel: logChannel,
		Db:                 db,
	}

	var projectNames []string
	for _, project := range config.Projects {
		projectNames = append(projectNames, project.Identifier)
	}
	runId, err := persistence.CreateNewRun(db, projectNames)
	if err != nil {
		log.Fatal(err)
	}

	service := run.CreateService(runId, dj)
	go service.Execute()
	go terminalui.TerminalUi(config, prgs, progressChannel, logChannel)
}

func test(progressChannel chan []progress.ProjectProgress, logChannel chan string) {
	for i := range 100 {
		time.Sleep(1 * time.Second)
		logChannel <- "Test " + strconv.Itoa(i)
	}
}
