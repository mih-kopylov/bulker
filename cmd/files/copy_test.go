package files

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testCopyResult struct {
	Repo   string `json:"repo"`
	Source string `json:"source"`
	Target string `json:"target"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func TestCopy(t *testing.T) {
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
	err = os.WriteFile(tests.Path("repo", "file.md"), []byte("hi"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCopyCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo --source file.md --target file2.md")
	assert.NoError(t, err)
	assert.Equal(t, "copy", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testCopyResult{
				{
					Repo:   "repo",
					Status: "copied",
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

func TestCopy_TargetFileAlreadyExists(t *testing.T) {
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
	err = os.WriteFile(tests.Path("repo", "file.md"), []byte("hi"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(tests.Path("repo", "file2.md"), []byte("old"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCopyCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo --source file.md --target file2.md")
	assert.NoError(t, err)
	assert.Equal(t, "copy", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testCopyResult{
				{
					Repo:   "repo",
					Status: "failed",
					Source: tests.Path("repo", "file.md"),
					Target: tests.Path("repo", "file2.md"),
					Error:  "target already exists",
				},
			},
		), output,
	)

	file2Content, err := os.ReadFile(tests.Path("repo", "file2.md"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("old"), file2Content)
}

func TestCopy_Force(t *testing.T) {
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
	err = os.WriteFile(tests.Path("repo", "file.md"), []byte("hi"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(tests.Path("repo", "file2.md"), []byte("old"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCopyCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo --source file.md --target file2.md --force")
	assert.NoError(t, err)
	assert.Equal(t, "copy", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testCopyResult{
				{
					Repo:   "repo",
					Status: "copied",
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

func TestCopy_SourceRequired(t *testing.T) {
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

	command := CreateCopyCommand(sh)
	_, _, err = tests.ExecuteCommand(command, "-n repo --target file2.md")
	if assert.Error(t, err) {
		assert.Equal(t, "required flag(s) \"source\" not set", err.Error())
	}
}

func TestCopy_TargetRequired(t *testing.T) {
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

	command := CreateCopyCommand(sh)
	_, _, err = tests.ExecuteCommand(command, "-n repo --source file.md")
	if assert.Error(t, err) {
		assert.Equal(t, "required flag(s) \"target\" not set", err.Error())
	}
}
