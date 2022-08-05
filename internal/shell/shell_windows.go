//go:build windows

package shell

import (
	"fmt"
	"os/exec"
	"syscall"
)

// RunCommand Runs a shell command in commandRootDirectory. If the commandRootDirectory is empty,
// runs the command from the current working directory.
// It returns combined stdout and stderr content, as it's visible in console
func RunCommand(commandRootDirectory string, command string, arguments ...string) (string, error) {
	cmd := exec.Command(command, arguments...)
	cmd.Dir = commandRootDirectory
	// this makes the child process ignore the SIGTERM for the bulker
	// so that bulker waits till the child command successfully completes and only then terminates
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil

}
