package tests

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/shell"
	"path/filepath"
)

type dynamicMockShell struct {
	handler func(repoName string, command string, arguments []string) (string, error)
}

func (m *dynamicMockShell) RunCommand(commandRootDirectory string, command string, arguments ...string) (
	string, error,
) {
	repoName := filepath.Base(commandRootDirectory)
	return m.handler(repoName, command, arguments)
}

// MockShellFunc returns a mock of a Shell interface that returns a conditional result based on command and arguments
func MockShellFunc(handler func(repoName string, command string, arguments []string) (string, error)) shell.Shell {
	return &dynamicMockShell{handler: handler}
}

type MockResult struct {
	Output string
	Error  error
}

func MockShellMap(handlerValues map[string]MockResult) shell.Shell {
	return MockShellFunc(
		func(repoName string, command string, arguments []string) (string, error) {
			commandLine := ShellCommandToString(command, arguments)
			mockResult, found := handlerValues[commandLine]
			if found {
				return mockResult.Output, mockResult.Error
			}
			return "", fmt.Errorf("shell not mocked: %v %v", repoName, commandLine)
		},
	)
}
