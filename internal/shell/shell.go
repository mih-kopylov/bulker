package shell

import (
	"fmt"
	"os/exec"
)

// RunCommand Runs a shell command in commandRootDirectory. If the commandRootDirectory is empty,
// runs the command from the current working directory.
// It returns combined stdout and stderr content, as it's visible in console
func RunCommand(commandRootDirectory string, command string, arguments ...string) (string, error) {
	cmd := exec.Command(command, arguments...)
	cmd.Dir = commandRootDirectory

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil

}
