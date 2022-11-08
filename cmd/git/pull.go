package git

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreatePullCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}

	var result = &cobra.Command{
		Use:   "pull",
		Short: "Pull changes from remote",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				gitService := gitops.NewGitService(sh)
				err := gitService.Pull(runContext.Repo)
				if err != nil {
					return nil, err
				}

				return "pulled successfully", nil
			},
		),
	}

	filter.AddCommandFlags(result)

	return result
}
