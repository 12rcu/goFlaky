package execution

import (
	"goFlaky/core"
	"os"
	"os/exec"
	"strings"
)

func OsCommand(command core.TestExecutionCommand, dir string, logs chan string, terminalLog chan string) {
	terminalLog <- "Executing " + command.Command + " " + strings.Join(command.Args, " ") + " In directory " + dir
	err := os.Chmod(dir+"/"+command.Command, 0777)
	if err != nil {
		terminalLog <- "[ERROR] while changing file perms " + err.Error()
	}
	cmd := exec.Command("./"+command.Command, command.Args...)
	cmd.Dir = dir
	cmd.Path = "/bin/sh"
	stdout, err := cmd.Output()
	if err != nil {
		terminalLog <- "[ERROR] " + err.Error()
	}
	terminalLog <- string(stdout)
	terminalLog <- "Finished executing " + command.Command + " " + strings.Join(command.Args, " ") + " In directory " + dir
}
