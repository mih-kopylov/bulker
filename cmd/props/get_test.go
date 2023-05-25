package props

import (
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

func TestGet(t *testing.T) {
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
	err = os.WriteFile(tests.Path("repo", "file.yaml"), []byte(`key: value`), os.ModePerm)
	assert.NoError(t, err)

	command := CreateGetCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -f file.yaml -p $.key")
	assert.NoError(t, err)
	assert.Equal(t, "get", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:   "repo",
					Result: "file.yaml: value",
				},
			},
		), output,
	)
}

func TestGet_KeyNotFound(t *testing.T) {
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
	err = os.WriteFile(tests.Path("repo", "file.yaml"), []byte(`key: value`), os.ModePerm)
	assert.NoError(t, err)

	command := CreateGetCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -f file.yaml -p $.key1")
	assert.NoError(t, err)
	assert.Equal(t, "get", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{},
		), output,
	)
}

func TestGet_FileNotFound(t *testing.T) {
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

	command := CreateGetCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -f file.yaml -p $.key")
	assert.NoError(t, err)
	assert.Equal(t, "get", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{},
		), output,
	)
}

func TestGet_RequiredFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    string
		message string
	}{
		{
			name:    "files",
			args:    "-p $.key",
			message: "required flag(s) \"files\" not set",
		},
		{
			name:    "path",
			args:    "-f file.yaml",
			message: "required flag(s) \"path\" not set",
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

				command := CreateGetCommand(sh)
				_, _, err = tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}

}
