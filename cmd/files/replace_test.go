package files

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testReplaceResult struct {
	Repo   string `json:"repo"`
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func TestReplace(t *testing.T) {
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
	err = os.WriteFile(tests.Path("repo", "file.md"), []byte("hi there"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(tests.Path("repo", "file2.md"), []byte("another hi there"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateReplaceCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -f *.md -c hi -r hello")
	assert.NoError(t, err)
	assert.Equal(t, "replace", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testReplaceResult{
				{
					Repo:   "repo",
					Result: "file.md :: 1\nfile2.md :: 1",
				},
			},
		), output,
	)

	fileContent, err := os.ReadFile(tests.Path("repo", "file.md"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello there"), fileContent)

	file2Content, err := os.ReadFile(tests.Path("repo", "file2.md"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("another hello there"), file2Content)
}

func TestReplace_RequiredFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    string
		message string
	}{
		{
			name:    "files",
			args:    "-n repo -c abc -r cba",
			message: "required flag(s) \"files\" not set",
		},
		{
			name:    "contains",
			args:    "-n repo -f *.md -r cba",
			message: "required flag(s) \"contains\" not set",
		},
		{
			name:    "replacement",
			args:    "-n repo -f *.md -c abc",
			message: "required flag(s) \"replacement\" not set",
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

				command := CreateReplaceCommand(sh)
				_, _, err = tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
