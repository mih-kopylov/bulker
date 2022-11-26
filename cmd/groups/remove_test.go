package groups

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemove(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{"repo"}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-g 1")
	if assert.NoError(t, err) {
		assert.Equal(t, "remove", c.Name())
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
		_, err := sets.GetGroup("1")
		if assert.Error(t, err) {
			assert.Equal(t, "group is not found", err.Error())
		}
	}
}

func TestRemove_NotExistingGroup(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{"repo"}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateRemoveCommand(sh)
	_, _, err := tests.ExecuteCommand(command, "-g 2")
	if assert.Error(t, err) {
		assert.Equal(t, "group is not found", err.Error())
	}
}

func TestRemove_RequiredFlags(t *testing.T) {
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

				command := CreateRemoveCommand(sh)
				_, _, err := tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
