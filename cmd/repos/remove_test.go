package repos

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testRemoveResult struct {
	Repo   string `json:"repo"`
	Result string `json:"result"`
}

func TestRemove(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com/tenant/repo.git"},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulker(t, sh, repos)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo")
	if assert.NoError(t, err) {
		assert.Equal(t, "remove", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testRemoveResult{
					{
						Repo:   "repo",
						Result: "removed",
					},
				},
			), output,
		)
	}

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	if assert.NoError(t, err) {
		_, err := sets.GetRepo("repo")
		if assert.Error(t, err) {
			assert.Equal(t, "repository is not found", err.Error())
		}
	}
}

func TestRemove_RepoNotFound(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com/tenant/repo.git"},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulker(t, sh, repos)

	command := CreateRemoveCommand(sh)
	_, _, err := tests.ExecuteCommand(command, "-n another")
	if assert.Error(t, err) {
		assert.Equal(t, "repository is not found", err.Error())
	}

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	if assert.NoError(t, err) {
		repo, err := sets.GetRepo("repo")
		if assert.NoError(t, err) {
			assert.Equal(t, "repo", repo.Name)
		}
	}
}

func TestRemove_RequiredFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    string
		message string
	}{
		{
			name:    "name",
			args:    "",
			message: "required flag(s) \"name\" not set",
		},
	}
	for _, test := range cases {
		t.Run(
			test.name, func(t *testing.T) {
				var repos []settings.Repo
				sh := tests.MockShellEmpty()
				tests.PrepareBulker(t, sh, repos)

				command := CreateRemoveCommand(sh)
				_, _, err := tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
