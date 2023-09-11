package repos

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testExportResult struct {
	Repo   string `json:"repo"`
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func TestExport_ToEmptyRepo(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := &shell.NativeShell{}
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(tests.Path("repo"), os.ModePerm)
	assert.NoError(t, err)
	bareGitRepo, err := tests.CreateBareGitRepository("likeRemoteGitRepo")
	assert.NoError(t, err)

	command := CreateExportCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-r "+bareGitRepo)
	if assert.NoError(t, err) {
		assert.Equal(t, "export", c.Name())
		assert.JSONEq(
			t, tests.ToJsonString(
				[]testExportResult{
					{
						Repo:   "repo",
						Result: "exported",
					},
				},
			), output,
		)

		reposFileContent, err := tests.GetFileContent(bareGitRepo, "repos.yaml")
		if assert.NoError(t, err) {
			assert.YAMLEq(
				t, `
version: 1
data:
    repos:
        repo:
            url: https://example.com
            tags: []
`, reposFileContent,
			)
		}
	}
}
