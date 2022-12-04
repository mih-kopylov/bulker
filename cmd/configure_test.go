package cmd

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

type testConfigureResult struct {
	Entity string `json:"entity"`
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func TestConfigure(t *testing.T) {
	sh := tests.MockShellEmpty()
	testDirectory := t.TempDir()
	testBulkerConfigFile := filepath.Join(testDirectory, "test_bulker.yaml")
	conf := &config.Config{
		ReposDirectory: testDirectory,
		MaxWorkers:     12,
		Output:         config.JsonOutputFormat,
	}
	err := config.WriteConfig(conf, testBulkerConfigFile)
	assert.NoError(t, err)
	viper.SetConfigFile(testBulkerConfigFile)
	err = viper.ReadInConfig()
	assert.NoError(t, err)
	assert.Equal(t, testDirectory, config.ReadConfig().ReposDirectory)
	assert.Equal(t, 12, config.ReadConfig().MaxWorkers)

	command := CreateRootCommand("", sh)
	c, output, err := tests.ExecuteCommand(
		command, fmt.Sprintf("configure --save %v --max-workers 13", testBulkerConfigFile),
	)
	if assert.NoError(t, err) {
		assert.Equal(t, "configure", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testConfigureResult{
					{
						Entity: "configuration",
						Result: "saved",
					},
				},
			), output,
		)

		err = viper.ReadInConfig()
		assert.NoError(t, err)

		assert.Equal(t, 13, config.ReadConfig().MaxWorkers)
	}
}
