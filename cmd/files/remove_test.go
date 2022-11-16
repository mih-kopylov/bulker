package files

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testRemoveResult struct {
	Repo    string `json:"repo"`
	Removed string `json:"removed"`
	Error   string `json:"error,omitempty"`
}

func TestRemove(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(tests.Path("repo", "file.md"), []byte("hi"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -f file.md")
	assert.NoError(t, err)
	assert.Equal(t, "remove", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testRemoveResult{
				{
					Repo:    "repo",
					Removed: fmt.Sprintf("%v: removed", tests.Path("repo", "file.md")),
				},
			},
		), output,
	)

	exists, err := utils.Exists(tests.Path("repo", "file.md"))
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestRemove_FileNotFound(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(tests.Path("repo", "file1.md"), []byte("hi"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -f file.md")
	assert.NoError(t, err)
	assert.Equal(t, "remove", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testRemoveResult{},
		), output,
	)

	exists, err := utils.Exists(tests.Path("repo", "file1.md"))
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRemove_Doublestar(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellEmpty()
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)
	err = os.MkdirAll(tests.Path("repo", "level1", "level2"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(tests.Path("repo", "level1", "level2", "file.md"), []byte("hi"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -f **/file.md")
	assert.NoError(t, err)
	assert.Equal(t, "remove", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testRemoveResult{
				{
					Repo:    "repo",
					Removed: fmt.Sprintf("%v: removed", tests.Path("repo", "level1", "level2", "file.md")),
				},
			},
		), output,
	)

	exists, err := utils.Exists(tests.Path("repo", "level1", "level2", "file.md"))
	assert.NoError(t, err)
	assert.False(t, exists)
}
