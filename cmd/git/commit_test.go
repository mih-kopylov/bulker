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

func TestCommit(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git add **":            {Output: "OK"},
			"git commit -m message": {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCommitCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -m message")
	assert.NoError(t, err)
	assert.Equal(t, "commit", c.Name())
	assert.Equal(t, output, "repo: result=committed\n")
}

func TestCommit_Pattern(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git add *.md":          {Output: "OK"},
			"git commit -m message": {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCommitCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -m message -p *.md")
	assert.NoError(t, err)
	assert.Equal(t, "commit", c.Name())
	assert.Equal(t, output, "repo: result=committed\n")
}
