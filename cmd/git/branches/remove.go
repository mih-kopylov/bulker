package branches

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
	"strings"
)

func CreateRemoveCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}
	var flags struct {
		name string
		mode config.GitMode
	}

	var result = &cobra.Command{
		Use:   "remove",
		Short: "Remove a branch",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				gitService := gitops.NewGitService(sh)
				branches, err := gitService.GetBranches(runContext.Repo, flags.mode, flags.name)
				if err != nil {
					return nil, err
				}

				var removeResult []string
				for _, branch := range branches {
					err := gitService.RemoveBranch(runContext.Repo, branch)
					if err != nil {
						removeResult = append(removeResult, fmt.Sprintf("%v: %v", branch.Short(), err.Error()))
					} else {
						removeResult = append(removeResult, fmt.Sprintf("%v: removed", branch.Short()))
					}
				}

				if len(removeResult) == 0 {
					return nil, nil
				}

				return strings.Join(removeResult, "\n"), nil

			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(&flags.name, "branch", "b", "", "Name of the branch to remove")
	utils.MarkFlagRequiredOrFail(result.Flags(), "branch")

	config.AddGitModeFlag(&flags.mode, result.Flags())
	return result
}
