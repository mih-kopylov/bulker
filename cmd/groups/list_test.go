package groups

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{"repo"}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateListCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "")
	if assert.NoError(t, err) {
		assert.Equal(t, "list", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testGroupResult{
					{
						Group: "1",
					},
				},
			), output,
		)
	}
}
