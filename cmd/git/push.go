package git

import (
	"context"
	"errors"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreatePushCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}

	var flags struct {
		force       bool
		allBranches bool
		branch      string
	}

	var result = &cobra.Command{
		Use:   "push",
		Short: "Push branches to remote",
		Long: `Push branches to remote.
If -b <branchName> is defined, pushes the only branch. Otherwise pushes all branches`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if flags.branch == "" && !flags.allBranches {
				return errors.New("either 'branch' or 'all' flags should be set")
			}
			if flags.allBranches && flags.force {
				return errors.New("only one branch is allowed to be force pushed")
			}
			return nil
		},
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				gitService := gitops.NewGitService(sh)
				err := gitService.Push(runContext.Repo, flags.branch, flags.allBranches, flags.force)
				if err != nil {
					return nil, err
				}

				return "Pushed", nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(&flags.branch, "branch", "b", "", "Name of the branch to push")

	result.Flags().BoolVarP(&flags.allBranches, "all", "a", false, "Push all branches if defined")

	result.MarkFlagsMutuallyExclusive("branch", "all")

	result.Flags().BoolVarP(&flags.force, "force", "f", false, "Use force push if defined")

	return result
}
