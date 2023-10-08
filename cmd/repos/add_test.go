package repos

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testAddResult struct {
	Repo   string `json:"repo"`
	Result string `json:"result"`
}

func TestAdd(t *testing.T) {
	var repos []settings.Repo
	sh := tests.MockShellEmpty()
	tests.PrepareBulker(t, sh, repos)

	command := CreateAddCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "--url https://example.com/tenant/repo.git --tags one")
	if assert.NoError(t, err) {
		assert.Equal(t, "add", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testAddResult{
					{
						Repo:   "repo",
						Result: "added",
					},
				},
			), output,
		)
	}

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	if assert.NoError(t, err) {
		repo, err := sets.GetRepo("repo")
		if assert.NoError(t, err) {
			assert.Equal(t, "repo", repo.Name)
			assert.Equal(t, "https://example.com/tenant/repo.git", repo.Url)
			assert.Equal(t, []string{"one"}, repo.Tags)
		}
	}
}

func TestAdd_RepoAlreadyExists(t *testing.T) {
	repos := []settings.Repo{
		{Name: "repo", Url: "https://example.com/tenant/repo.git"},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulker(t, sh, repos)

	command := CreateAddCommand(sh)
	_, _, err := tests.ExecuteCommand(command, "--url https://example.com/tenant/repo.git")
	if assert.Error(t, err) {
		assert.Equal(t, "repository already exists", err.Error())
	}
}

func TestAdd_NameOverrides(t *testing.T) {
	var repos []settings.Repo
	sh := tests.MockShellEmpty()
	tests.PrepareBulker(t, sh, repos)

	command := CreateAddCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "--url https://example.com/tenant/repo.git --name another")
	if assert.NoError(t, err) {
		assert.Equal(t, "add", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testAddResult{
					{
						Repo:   "another",
						Result: "added",
					},
				},
			), output,
		)
	}

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	if assert.NoError(t, err) {
		repo, err := sets.GetRepo("another")
		if assert.NoError(t, err) {
			assert.Equal(t, "another", repo.Name)
			assert.Equal(t, "https://example.com/tenant/repo.git", repo.Url)
			assert.Equal(t, []string{}, repo.Tags)
		}

		_, err = sets.GetRepo("repo")
		if assert.Error(t, err) {
			assert.Equal(t, "repository is not found", err.Error())
		}
	}
}

func TestAdd_MultipleTags(t *testing.T) {
	var repos []settings.Repo
	sh := tests.MockShellEmpty()
	tests.PrepareBulker(t, sh, repos)

	command := CreateAddCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "--url https://example.com/tenant/repo.git --tags one --tags two")
	if assert.NoError(t, err) {
		assert.Equal(t, "add", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testAddResult{
					{
						Repo:   "repo",
						Result: "added",
					},
				},
			), output,
		)
	}

	manager := settings.NewManager(config.ReadConfig(), sh)
	sets, err := manager.Read()
	if assert.NoError(t, err) {
		repo, err := sets.GetRepo("repo")
		if assert.NoError(t, err) {
			assert.Equal(t, "repo", repo.Name)
			assert.Equal(t, "https://example.com/tenant/repo.git", repo.Url)
			assert.Equal(t, []string{"one", "two"}, repo.Tags)
		}
	}
}

func TestAdd_RequiredFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    string
		message string
	}{
		{
			name:    "url",
			args:    "",
			message: "required flag(s) \"url\" not set",
		},
	}
	for _, test := range cases {
		t.Run(
			test.name, func(t *testing.T) {
				var repos []settings.Repo
				sh := tests.MockShellEmpty()
				tests.PrepareBulker(t, sh, repos)

				command := CreateAddCommand(sh)
				_, _, err := tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
