package git

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
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
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, _, err := tests.ExecuteCommand(command, "-n repo")
	assert.Error(t, err, "either 'branch' or 'all' flags should be set")
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
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b branch-name")
	assert.NoError(t, err)
	assert.Equal(t, "push", c.Name())
	assert.Equal(t, output, "repo: result=Pushed\n")
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
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo --all")
	assert.NoError(t, err)
	assert.Equal(t, "push", c.Name())
	assert.Equal(t, output, "repo: result=Pushed\n")
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
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b my-branch -f")
	assert.NoError(t, err)
	assert.Equal(t, "push", c.Name())
	assert.Equal(t, output, "repo: result=Pushed\n")
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
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreatePushCommand(sh)
	c, _, err := tests.ExecuteCommand(command, "-n repo --all -f")
	assert.Error(t, err, "only one branch is allowed to be force pushed")
	assert.Equal(t, "push", c.Name())
}
