package branches

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateRemoveCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}
	var flags struct {
		name string
		mode gitops.GitMode
	}

	var result = &cobra.Command{
		Use:   "remove",
		Short: "Remove a branch",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				gitService := gitops.NewGitService(sh)
				removeResult, err := gitService.RemoveBranch(runContext.Repo, flags.name, flags.mode)
				if err != nil {
					return nil, err
				}

				return removeResult, nil

			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(&flags.name, "branch", "b", "", "Name of the branch to remove")
	utils.MarkFlagRequiredOrFail(result.Flags(), "branch")

	flags.mode = gitops.GitModeAll
	result.Flags().VarP(
		&flags.mode, "mode", "m", fmt.Sprintf(
			"Type of branches to process. "+
				"Available types are: %s, %s, %s", gitops.GitModeAll, gitops.GitModeLocal, gitops.GitModeRemote,
		),
	)

	return result
}
