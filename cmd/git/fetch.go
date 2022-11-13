package git

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateFetchCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}

	var result = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch changes from remote",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				gitService := gitops.NewGitService(sh)
				err := gitService.Fetch(runContext.Repo)
				if err != nil {
					return nil, err
				}

				return "Fetched", nil
			},
		),
	}

	filter.AddCommandFlags(result)

	return result
}
