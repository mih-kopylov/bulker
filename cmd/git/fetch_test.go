package git

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git fetch --prune": {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateFetchCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo")
	assert.NoError(t, err)
	assert.Equal(t, "fetch", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "fetched",
				},
			},
		), output,
	)
}
