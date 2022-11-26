package groups

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{"repo"}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateGetCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-g 1")
	if assert.NoError(t, err) {
		assert.Equal(t, "get", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testResult{
					{
						Repo: "repo",
					},
				},
			), output,
		)
	}

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	if assert.NoError(t, err) {
		group, err := sets.GetGroup("1")
		if assert.NoError(t, err) {
			assert.Equal(t, []string{"repo"}, group.Repos)
		}
	}
}

func TestGet_RequiredFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    string
		message string
	}{
		{
			name:    "group",
			args:    "",
			message: "required flag(s) \"group\" not set",
		},
	}
	for _, test := range cases {
		t.Run(
			test.name, func(t *testing.T) {
				repos := []settings.Repo{
					{
						Name: "repo",
						Url:  "https://example.com",
					},
				}
				groups := []settings.Group{
					{Name: "1", Repos: []string{"repo"}},
				}
				sh := tests.MockShellEmpty()
				tests.PrepareBulkerWithGroups(t, sh, repos, groups)

				command := CreateGetCommand(sh)
				_, _, err := tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
