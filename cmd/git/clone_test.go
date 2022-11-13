package git

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestClone(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git status":                      {Output: "OK"},
			"git clone https://example.com .": {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.Equal(t, "repo: result=cloned\n", output)
}

func TestClone_EmptyDirectoryExists(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git status":                      {Output: "OK"},
			"git clone https://example.com .": {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.Equal(t, "repo: result=cloned\n", output)
}

func TestClone_NotEmptyDirectoryExists(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git status": {Output: "err", Error: fmt.Errorf("not a repository")},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(viper.GetString("reposDirectory"), "repo", "file"), []byte("data"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.Equal(t, "repo: failed to clone: failed to get git status: err, not a repository\n", output)
}

func TestClone_Recreate(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git status": {
				Output: "On branch custom\nYour branch is up to date with 'origin/custom'.",
			},
			"git clone https://example.com .": {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(viper.GetString("reposDirectory"), "repo", "file"), []byte("data"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo --recreate")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.Equal(t, "repo: result=re-cloned\n", output)
}
