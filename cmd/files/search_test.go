package files

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testSearchResult struct {
	Repo   string `json:"repo"`
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func TestSearch(t *testing.T) {
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
	err = os.WriteFile(tests.Path("repo", "file.md"), []byte("first\nhello\nthird"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateSearchCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -f *.md -c hello -b 1 -a 1")
	assert.NoError(t, err)
	assert.Equal(t, "search", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testSearchResult{
				{
					Repo: "repo",
					Result: `file.md:
    - |-
      first
      hello
      third`,
				},
			},
		), output,
	)
}

func TestSearch_RequiredFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    string
		message string
	}{
		{
			name:    "files",
			args:    "-n repo",
			message: "required flag(s) \"files\" not set",
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

				command := CreateSearchCommand(sh)
				_, _, err = tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
