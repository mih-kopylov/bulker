package git

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPush(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, _, err := tests.ExecuteCommand(command, "-n repo")
	if assert.Error(t, err) {
		assert.Equal(t, "either 'branch' or 'all' flags should be set", err.Error())
	}
	assert.Equal(t, "push", c.Name())
}

func TestPush_Branch(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git push --set-upstream origin branch-name": {Output: "OK"},
			"git remote": {Output: "origin"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b branch-name")
	assert.NoError(t, err)
	assert.Equal(t, "push", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "pushed",
				},
			},
		), output,
	)
}

func TestPush_All(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git push --set-upstream origin --all": {Output: "OK"},
			"git remote":                           {Output: "origin"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo --all")
	assert.NoError(t, err)
	assert.Equal(t, "push", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "pushed",
				},
			},
		), output,
	)
}

func TestPush_Branch_Force(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git push --set-upstream origin my-branch --force": {Output: "OK"},
			"git remote": {Output: "origin"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b my-branch -f")
	assert.NoError(t, err)
	assert.Equal(t, "push", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "pushed",
				},
			},
		), output,
	)
}

func TestPush_All_Force(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, _, err := tests.ExecuteCommand(command, "-n repo --all -f")
	if assert.Error(t, err) {
		assert.Equal(t, "only one branch is allowed to be force pushed", err.Error())
	}
	assert.Equal(t, "push", c.Name())
}
