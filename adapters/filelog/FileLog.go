package filelog

import (
	"goFlaky/core"
	"log"
	"os"
	"strconv"
)

func FileLog(config *core.Config, logChannel chan string, runId int) {
	path := config.BaseDir + "/" + config.LogDir + "/" + strconv.Itoa(runId) + "log.txt"
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Cannot open log file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Cannot close log file: %v", err)
		}
	}(f)

	for logEntry := range logChannel {
		_, err := f.WriteString(logEntry + "\n")
		if err != nil {
			log.Printf("Cannot write log file: %v", err)
		}
	}
}
