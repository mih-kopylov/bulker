package shell

type Shell interface {
	RunCommand(commandRootDirectory string, command string, arguments ...string) (string, error)
}

type NativeShell struct {
}
