package tests

import (
	"bytes"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func ExecuteCommand(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)

	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func PrepareBulker(t *testing.T, sh shell.Shell, repos []settings.Repo) {
	testDirectory := t.TempDir()
	viper.Set("reposDirectory", testDirectory)
	viper.Set("settingsFileName", filepath.Join(testDirectory, "bulker_test_settings.yaml"))
	viper.Set("runMode", "seq")
	viper.Set("noProgress", "true")
	viper.Set("output", "line")

	conf := config.ReadConfig()
	manager := settings.NewManager(conf, sh)
	s := &settings.Settings{
		Repos:  repos,
		Groups: nil,
	}
	err := manager.Write(s)
	assert.NoError(t, err)
}

func ShellCommandEquals(actualCommand string, actualArguments []string, expected string) bool {
	buffer := bytes.Buffer{}
	buffer.WriteString(actualCommand)
	buffer.WriteString(" ")
	buffer.WriteString(strings.Join(actualArguments, " "))
	return buffer.String() == expected
}
