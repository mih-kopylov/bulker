package branches

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testResult struct {
	Repo     string `json:"repo"`
	Error    string `json:"error,omitempty"`
	Status   string `json:"status"`
	Checkout string `json:"checkout"`
	Ref      string `json:"ref"`
}

func TestCheckout(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git branch -a --format=%(refname)": {Output: "refs/heads/master\nrefs/heads/br"},
			"git checkout br":                   {Output: "Switched to branch 'br'"},
			"git status":                        {Output: "On branch br\nnothing to commit, working tree clean"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCheckoutCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b br")
	assert.NoError(t, err)
	assert.Equal(t, "checkout", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:     "repo",
					Status:   "clean",
					Ref:      "br",
					Checkout: "success",
				},
			},
		), output,
	)
}

func TestCheckout_FromDetachedHead(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git branch -a --format=%(refname)": {
				Output: "(HEAD detached at 0123456)\nrefs/heads/master\nrefs/heads/br",
			},
			"git checkout br": {Output: "Switched to branch 'br'"},
			"git status":      {Output: "On branch br\nnothing to commit, working tree clean"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCheckoutCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b br")
	assert.NoError(t, err)
	assert.Equal(t, "checkout", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:     "repo",
					Status:   "clean",
					Ref:      "br",
					Checkout: "success",
				},
			},
		), output,
	)
}

func TestCheckout_Discard(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := tests.MockShellMap(
		map[string]tests.MockResult{
			"git reset --hard HEAD":             {Output: "OK"},
			"git branch -a --format=%(refname)": {Output: "refs/heads/master\nrefs/heads/br"},
			"git checkout br":                   {Output: "Switched to branch 'br'"},
			"git status":                        {Output: "On branch br\nnothing to commit, working tree clean"},
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCheckoutCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n repo -b br -d")
	assert.NoError(t, err)
	assert.Equal(t, "checkout", c.Name())
	assert.JSONEq(
		t, tests.ToJsonString(
			[]testResult{
				{
					Repo:     "repo",
					Status:   "clean",
					Ref:      "br",
					Checkout: "success",
				},
			},
		), output,
	)
}

func TestCheckout_RequiredFlags(t *testing.T) {
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

				command := CreateCheckoutCommand(sh)
				_, _, err = tests.ExecuteCommand(command, test.args)
				if assert.Error(t, err) {
					assert.Equal(t, test.message, err.Error())
				}
			},
		)
	}
}
