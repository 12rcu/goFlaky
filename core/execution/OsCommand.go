package execution

import (
	"os/exec"
)

func OsCommand(command string, dir string, logs chan string, terminalLog chan string) {
	terminalLog <- "Executing " + command + " In directory " + dir
	cmd := exec.Command(command)
	cmd.Dir = dir
	stdout, err := cmd.Output()
	if err != nil {
		terminalLog <- err.Error()
	}
	logs <- string(stdout)

	if err := cmd.Run(); err != nil {
		terminalLog <- err.Error()
	}
	terminalLog <- "Finished executing " + command + " In directory " + dir
}
