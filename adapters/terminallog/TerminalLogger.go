package terminallog

import (
	"goFlaky/core/progress"
	"log"
	"sync"
)

func TerminalLogger(
	logChannel chan string,
	progressChannel chan []progress.ProjectProgress,
	waitGroup *sync.WaitGroup,
) {
	go consumeProgress(progressChannel)
	receiveLogs(logChannel)
	waitGroup.Done()
}

func consumeProgress(progressChannel chan []progress.ProjectProgress) {
	for _ = range progressChannel {
	}
}

func receiveLogs(logChannel chan string) {
	for prg := range logChannel {
		log.Printf("[LOG] %s\n", prg)
	}
}
