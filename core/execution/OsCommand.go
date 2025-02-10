package execution

import (
	"bytes"
	"goFlaky/core"
	"os/exec"
	"strings"
)

func OsCommand(command core.TestExecutionCommand, dir string, logs chan string, terminalLog chan string) {
	terminalLog <- "Executing " + command.Command + " " + strings.Join(command.Args, " ") + " In directory " + dir
	cmd := exec.Command(command.Command, command.Args...)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Dir = dir
	cmd.Path = command.Path

	err := cmd.Run()
	if err != nil {
		terminalLog <- "[ERROR] " + err.Error()
	}
	logs <- "====== Command Output ======"
	logs <- outb.String()
	logs <- "=> Command Output (ERR)"
	logs <- errb.String()
	logs <- "====== End ======"
	terminalLog <- "Finished executing " + command.Command + " " + strings.Join(command.Args, " ") + " In directory " + dir
}
