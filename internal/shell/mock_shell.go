package shell

import "path/filepath"

type staticMockShell struct {
	resultString string
	resultError  error
}

func (m *staticMockShell) RunCommand(_ string, _ string, _ ...string) (string, error) {
	return m.resultString, m.resultError
}

// MockSuccess returns a mock of a Shell interface that succeeds with a provided value
func MockSuccess(result string) Shell {
	return &staticMockShell{result, nil}
}

// MockError returns a mock of a Shell interface that fails with a provided error
func MockError(result error) Shell {
	return &staticMockShell{"", result}
}

type dynamicMockShell struct {
	handler func(repoName string, command string, arguments []string) (string, error)
}

func (m *dynamicMockShell) RunCommand(commandRootDirectory string, command string, arguments ...string) (
	string, error,
) {
	repoName := filepath.Base(commandRootDirectory)
	return m.handler(repoName, command, arguments)
}

// MockShell returns a mock of a Shell interface that returns a conditional result based on command and arguments
func MockShell(handler func(repoName string, command string, arguments []string) (string, error)) Shell {
	return &dynamicMockShell{handler: handler}
}
