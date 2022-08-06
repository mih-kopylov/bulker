package git

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/spf13/cobra"
)

func CreatePushCommand() *cobra.Command {
	var filter = runner.Filter{}

	var flags struct {
		branch string
	}

	var result = &cobra.Command{
		Use:   "push",
		Short: "Push branches to remote",
		Long: `Push branches to remote.
If -b <branchName> is defined, pushes the only branch. Otherwise pushes all branches`,
		RunE: runner.NewDefaultRunner(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				err := gitops.Push(runContext.Fs, runContext.Repo, flags.branch)
				if err != nil {
					return nil, err
				}

				return "pushed successfully", nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(&flags.branch, "branch", "b", "", "Name of the branch to push")

	return result
}
