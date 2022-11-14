package git

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testResult struct {
	Repo   string `json:"repo"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

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
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "cloned",
				},
			},
		), output,
	)
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
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "cloned",
				},
			},
		), output,
	)
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
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(tests.Path("repo", "file"), []byte("data"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:  "repo",
					Error: "failed to clone: failed to get git status: err, not a repository",
				},
			},
		), output,
	)
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
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(tests.Path("repo", "file"), []byte("data"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo --recreate")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "re-cloned",
				},
			},
		), output,
	)
}
