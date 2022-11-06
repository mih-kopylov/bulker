package git

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/spf13/cobra"
)

func CreatePullCommand() *cobra.Command {
	var filter = runner.Filter{}

	var result = &cobra.Command{
		Use:   "pull",
		Short: "Pull changes from remote",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				err := gitops.Pull(runContext.Repo)
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
