package groups

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testResult struct {
	Repo   string `json:"repo"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func TestAppend(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateAppendCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-g 1 -n repo")
	assert.NoError(t, err)
	assert.Equal(t, "append", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "added",
				},
			},
		), output,
	)

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	assert.NoError(t, err)
	group, err := sets.GetGroup("1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"repo"}, group.Repos)
}

func TestAppend_GroupNotExists(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "2", Repos: []string{}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateAppendCommand(sh)
	_, _, err := tests.ExecuteCommand(command, "-g 1 -n repo")
	if assert.Error(t, err) {
		assert.Equal(t, "group is not found", err.Error())
	}

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	if assert.NoError(t, err) {
		group, err := sets.GetGroup("1")
		if assert.Error(t, err) {
			assert.Equal(t, "group is not found", err.Error())
		}

		group, err = sets.GetGroup("2")
		if assert.NoError(t, err) {
			assert.Equal(t, []string{}, group.Repos)
		}
	}
}

func TestAppend_RepoNotExists(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateAppendCommand(sh)
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

func TestAppend_NoRepoPassed(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateAppendCommand(sh)
	_, output, err := tests.ExecuteCommand(command, "-g 1")
	if assert.NoError(t, err) {
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testResult{},
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

func TestAppend_RequiredFlags(t *testing.T) {
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
					{Name: "1", Repos: []string{}},
				}
				sh := tests.MockShellEmpty()
				tests.PrepareBulkerWithGroups(t, sh, repos, groups)

				command := CreateAppendCommand(sh)
				_, _, err := tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
