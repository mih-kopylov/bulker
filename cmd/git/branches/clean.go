package branches

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateCleanCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}
	var flags struct {
		mode config.GitMode
	}

	var result = &cobra.Command{
		Use:   "clean",
		Short: "Remove all branches that are merged to default one",
		Long: `Remove all branches that are merged to default one.
First, it defines the default branch of the remote.
Then, it loops over the branches and removes the ones that don't have differences with the default one`,
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				gitService := gitops.NewGitService(sh)
				cleanResult, err := gitService.CleanBranches(runContext.Repo, flags.mode)
				if err != nil {
					return nil, err
				}

				if cleanResult == "" {
					return nil, nil
				}

				return cleanResult, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	config.AddGitModeFlag(&flags.mode, result.Flags())

	return result
}
