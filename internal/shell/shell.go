package shell

import (
	"bytes"
	"fmt"
	"os/exec"
)

// RunCommand Runs a shell command in commandRootDirectory. If the commandRootDirectory is empty,
// runs the command from the current working directory
func RunCommand(commandRootDirectory string, command string, arguments ...string) (string, error) {
	var errorOutput bytes.Buffer

	cmd := exec.Command(command, arguments...)
	cmd.Dir = commandRootDirectory
	cmd.Stderr = &errorOutput

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command failed: error output=%v, error=%w", errorOutput.String(), err)
	}

	return string(output), nil

}
