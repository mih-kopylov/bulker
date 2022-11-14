package git

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestPull(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git pull --prune": {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePullCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo")
	assert.NoError(t, err)
	assert.Equal(t, "pull", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "pulled",
				},
			},
		), output,
	)
}
