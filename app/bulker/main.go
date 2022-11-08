package main

import (
	"github.com/mih-kopylov/bulker/cmd"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/sirupsen/logrus"
)

var version string

func main() {
	rootCmd := cmd.CreateRootCommand(version, &shell.NativeShell{})
	err := rootCmd.Execute()
	if err != nil {
		logrus.Debugf("command failed: %v", err)
	}
}
