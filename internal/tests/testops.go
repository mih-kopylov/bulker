package tests

import (
	"bytes"
	"encoding/json"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"path/filepath"
	"strings"
	"testing"
)

func ExecuteCommand(root *cobra.Command, args string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)

	root.SetArgs(strings.Split(args, " "))

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func PrepareBulker(t *testing.T, sh shell.Shell, repos []settings.Repo) {
	testDirectory := t.TempDir()
	viper.Set("reposDirectory", testDirectory)
	viper.Set("settingsFileName", filepath.Join(testDirectory, "bulker_test_settings.yaml"))
	viper.Set("runMode", "seq")
	viper.Set("noProgress", "true")
	viper.Set("output", "json")

	conf := config.ReadConfig()
	manager := settings.NewManager(conf, sh)
	s := &settings.Settings{
		Repos:  repos,
		Groups: nil,
	}
	err := manager.Write(s)
	assert.NoError(t, err)
}

func ShellCommandToString(command string, arguments []string) string {
	buffer := bytes.Buffer{}
	buffer.WriteString(command)
	buffer.WriteString(" ")
	buffer.WriteString(strings.Join(arguments, " "))
	return buffer.String()
}

func ToJsonString(value any) string {
	marshalled, err := json.Marshal(value)
	if err != nil {
		logrus.Fatal(err)
	}
	return string(marshalled)
}

func Path(parts ...string) string {
	allParts := slices.Insert(parts, 0, viper.GetString("reposDirectory"))
	return filepath.Join(allParts...)
}
