package main

import (
	"goFlaky/adapters/filelog"
	"goFlaky/adapters/persistence"
	"goFlaky/adapters/terminallog"
	"goFlaky/adapters/terminalui"
	"goFlaky/core"
	"goFlaky/core/progress"
	"goFlaky/core/run"
	"log"
	"sync"
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
	fileLogChannel := make(chan string)

	dj := core.DependencyInjection{
		Config:             config,
		Progress:           prgs,
		ProgressChannel:    progressChannel,
		TerminalLogChannel: logChannel,
		Db:                 db,
		FileLogChannel:     fileLogChannel,
	}

	var projectNames []string
	for _, project := range config.Projects {
		projectNames = append(projectNames, project.Identifier)
	}
	runId, err := db.CreateNewRun(projectNames)
	if err != nil {
		log.Fatal(err)
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)

	go filelog.FileLog(config, fileLogChannel, runId)
	service := run.CreateService(runId, dj)
	go service.Execute(&waitGroup)

	if config.EnableUI {
		go terminalui.TerminalUi(prgs, progressChannel, logChannel, &waitGroup)
	} else {
		go terminallog.TerminalLogger(logChannel, progressChannel, &waitGroup)
	}

	waitGroup.Wait()
}
