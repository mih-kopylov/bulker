package groups

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, nil)

	command := CreateCreateCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-g 1 -n repo")
	if assert.NoError(t, err) {
		assert.Equal(t, "create", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testResult{
					{
						Repo:   "repo",
						Result: "created",
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

func TestCreate_RepoNotExists(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, nil)

	command := CreateCreateCommand(sh)
	_, output, err := tests.ExecuteCommand(command, "-g 1 -n repo1")
	if assert.NoError(t, err) {
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testResult{
					{
						Repo:  "repo1",
						Error: "repository is not supported",
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
			assert.Equal(t, []string{}, group.Repos)
		}
	}
}

func TestCreate_RequiredFlags(t *testing.T) {
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
				sh := tests.MockShellEmpty()
				tests.PrepareBulkerWithGroups(t, sh, repos, nil)

				command := CreateCreateCommand(sh)
				_, _, err := tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
