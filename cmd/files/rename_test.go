package files

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRename(t *testing.T) {
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

	command := CreateRenameCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo --source file.md --target file2.md")
	assert.NoError(t, err)
	assert.Equal(t, "rename", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testCopyResult{
				{
					Repo:   "repo",
					Status: "renamed",
					Source: tests.Path("repo", "file.md"),
					Target: tests.Path("repo", "file2.md"),
				},
			},
		), output,
	)

	exists, err := utils.Exists(tests.Path("repo", "file.md"))
	assert.NoError(t, err)
	assert.False(t, exists)

	file2Content, err := os.ReadFile(tests.Path("repo", "file2.md"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("hi"), file2Content)
}

func TestRename_SourceRequired(t *testing.T) {
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

	command := CreateRenameCommand(sh)
	_, _, err = tests.ExecuteCommand(command, "-n repo --target file2.md")
	if assert.Error(t, err) {
		assert.Equal(t, "required flag(s) \"source\" not set", err.Error())
	}
}

func TestRename_TargetRequired(t *testing.T) {
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

	command := CreateRenameCommand(sh)
	_, _, err = tests.ExecuteCommand(command, "-n repo --source file.md")
	if assert.Error(t, err) {
		assert.Equal(t, "required flag(s) \"target\" not set", err.Error())
	}
}

func TestRename_SourceNotFound(t *testing.T) {
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

	command := CreateRenameCommand(sh)
	_, output, err := tests.ExecuteCommand(command, "-n repo --source file.md --target file2.md")
	assert.NoError(t, err)
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testCopyResult{},
		), output,
	)

	exists, err := utils.Exists(tests.Path("repo", "file2.md"))
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestRename_TargetAlreadyExists(t *testing.T) {
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
	err = os.WriteFile(tests.Path("repo", "file2.md"), []byte("hello"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRenameCommand(sh)
	_, output, err := tests.ExecuteCommand(command, "-n repo --source file.md --target file2.md")
	assert.NoError(t, err)
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testCopyResult{
				{
					Repo:   "repo",
					Status: "failed",
					Error:  "target already exists",
					Source: tests.Path("repo", "file.md"),
					Target: tests.Path("repo", "file2.md"),
				},
			},
		), output,
	)

	file2Content, err := os.ReadFile(tests.Path("repo", "file2.md"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello"), file2Content)
}

func TestRename_Force(t *testing.T) {
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
	err = os.WriteFile(tests.Path("repo", "file2.md"), []byte("hello"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRenameCommand(sh)
	_, output, err := tests.ExecuteCommand(command, "-n repo --source file.md --target file2.md --force")
	assert.NoError(t, err)
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testCopyResult{
				{
					Repo:   "repo",
					Status: "renamed",
					Source: tests.Path("repo", "file.md"),
					Target: tests.Path("repo", "file2.md"),
				},
			},
		), output,
	)

	file2Content, err := os.ReadFile(tests.Path("repo", "file2.md"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("hi"), file2Content)
}
