package groups

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testGroupResult struct {
	Group  string `json:"group"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func TestClean(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{"repo"}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateCleanCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "")
	if assert.NoError(t, err) {
		assert.Equal(t, "clean", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testGroupResult{
					{
						Group:  "1",
						Result: "removed",
					},
				},
			), output,
		)
	}

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	if assert.NoError(t, err) {
		_, err = sets.GetGroup("1")
		if assert.Error(t, err) {
			assert.Equal(t, "group is not found", err.Error())
		}
	}
}
