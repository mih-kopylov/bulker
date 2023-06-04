package branches

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type removeTestResult struct {
	Repo   string `json:"repo"`
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
}

func TestRemove_Local(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git branch -a --format=%(refname)": {Output: "refs/heads/master\nrefs/heads/br"},
			"git branch -D br":                  {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b br -m all")
	assert.NoError(t, err)
	assert.Equal(t, "remove", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]removeTestResult{
				{
					Repo:   "repo",
					Result: "br: removed",
				},
			},
		), output,
	)
}

func TestRemove_Remote(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git branch -a --format=%(refname)": {Output: "refs/heads/master\nrefs/remotes/origin/br"},
			"git push origin --delete br":       {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b br -m all")
	assert.NoError(t, err)
	assert.Equal(t, "remove", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]removeTestResult{
				{
					Repo:   "repo",
					Result: "origin/br: removed",
				},
			},
		), output,
	)
}

func TestRemove_RemoteWhenOnlyLocalFound(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git branch -a --format=%(refname)": {Output: "refs/heads/master\nrefs/heads/br"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b br -m remote")
	assert.NoError(t, err)
	assert.Equal(t, "remove", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]removeTestResult{},
		), output,
	)
}

func TestRemove_LocalWhenOnlyRemoteFound(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git branch -a --format=%(refname)": {Output: "refs/heads/master\nrefs/remotes/origin/br"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b br -m local")
	assert.NoError(t, err)
	assert.Equal(t, "remove", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]removeTestResult{},
		), output,
	)
}

func TestRemove_NotFound(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git branch -a --format=%(refname)": {Output: "refs/heads/master\nrefs/heads/br1"},
			"git branch -D br":                  {Output: "OK"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b br")
	assert.NoError(t, err)
	assert.Equal(t, "remove", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]removeTestResult{},
		), output,
	)
}

func TestRemove_FailedToRemove(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git branch -a --format=%(refname)": {Output: "refs/heads/master\nrefs/heads/br\nrefs/remotes/origin/br"},
			"git branch -D br":                  {Output: "fail output", Error: fmt.Errorf("fail reason")},
			"git push origin --delete br":       {Output: "fail output", Error: fmt.Errorf("fail reason")},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateRemoveCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b br -m all")
	assert.NoError(t, err)
	assert.Equal(t, "remove", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]removeTestResult{
				{
					Repo: "repo",
					Result: `br: failed to remove local branch: fail output fail reason
origin/br: failed to remove remote branch: fail output fail reason`,
				},
			},
		), output,
	)
}

func TestRemove_RequiredFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    string
		message string
	}{
		{
			name:    "branch",
			args:    "-n repo",
			message: "required flag(s) \"branch\" not set",
		},
		{
			name:    "mode",
			args:    "-n repo -b br -m another",
			message: `invalid argument "another" for "-m, --mode" flag: must be one of 'all' 'local' 'remote'`,
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

				command := CreateRemoveCommand(sh)
				_, _, err = tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
