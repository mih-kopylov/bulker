package files

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/mih-kopylov/bulker/internal/utils"
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
	sh := tests.MockShellEmpty()
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
	sh := tests.MockShellEmpty()
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

func TestCopy_SourceNotFound(t *testing.T) {
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

	command := CreateCopyCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo --source file.md --target file2.md")
	assert.NoError(t, err)
	assert.Equal(t, "copy", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testCopyResult{},
		), output,
	)

	exists, err := utils.Exists(tests.Path("repo", "file2.md"))
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCopy_Force(t *testing.T) {
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

func TestCopy_RequiredFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    string
		message string
	}{
		{
			name:    "source",
			args:    "-n repo --target file2.md",
			message: "required flag(s) \"source\" not set",
		},
		{
			name:    "target",
			args:    "-n repo --source file.md",
			message: "required flag(s) \"target\" not set",
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
				sh := tests.MockShellEmpty()
				tests.PrepareBulker(t, sh, repos)
				err := os.Mkdir(tests.Path("repo"), os.ModePerm)
				assert.NoError(t, err)

				command := CreateCopyCommand(sh)
				_, _, err = tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
