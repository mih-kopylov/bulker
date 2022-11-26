package groups

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExclude(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{"repo"}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateExcludeCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-g 1 -n repo")
	if assert.NoError(t, err) {
		assert.Equal(t, "exclude", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testResult{
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
		group, err := sets.GetGroup("1")
		if assert.NoError(t, err) {
			assert.Equal(t, []string{}, group.Repos)
		}
	}
}

func TestExclude_GroupNotExists(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "2", Repos: []string{"repo"}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateExcludeCommand(sh)
	_, _, err := tests.ExecuteCommand(command, "-g 1 -n repo")
	if assert.Error(t, err) {
		assert.Equal(t, "group is not found", err.Error())
	}

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	if assert.NoError(t, err) {
		_, err := sets.GetGroup("1")
		if assert.Error(t, err) {
			assert.Equal(t, "group is not found", err.Error())
		}

		group, err := sets.GetGroup("2")
		if assert.NoError(t, err) {
			assert.Equal(t, []string{"repo"}, group.Repos)
		}
	}
}

func TestExclude_RepoNotExists(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateExcludeCommand(sh)
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

func TestExclude_RepoNotInGroup(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateExcludeCommand(sh)
	_, output, err := tests.ExecuteCommand(command, "-g 1 -n repo")
	if assert.NoError(t, err) {
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testResult{
					{
						Repo:   "repo",
						Result: "removing skipped",
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

func TestExclude_NoRepoPassed(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com"},
	}
	groups := []settings.Group{
		{Name: "1", Repos: []string{}},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulkerWithGroups(t, sh, repos, groups)

	command := CreateExcludeCommand(sh)
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

func TestExclude_RequiredFlags(t *testing.T) {
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

				command := CreateExcludeCommand(sh)
				_, _, err := tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
