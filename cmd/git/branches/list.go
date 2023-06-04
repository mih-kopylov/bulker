package branches

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
	"strings"
)

func CreateListCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}
	var flags struct {
		mode    config.GitMode
		pattern string
	}

	var result = &cobra.Command{
		Use:   "list",
		Short: "Prints a list of repository branches",
		Long: `Prints a list of repository branches.
If a repository doesn't have any branch matching pattern, the repository will be omitted in the result`,
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				gitService := gitops.NewGitService(sh)
				branches, err := gitService.GetBranches(runContext.Repo, flags.mode, flags.pattern)
				if err != nil {
					return nil, err
				}

				if len(branches) == 0 {
					return nil, nil
				}

				builder := strings.Builder{}
				for _, branch := range branches {
					builder.WriteString(branch.Short() + "\n")
				}

				return strings.TrimSpace(builder.String()), nil
			},
		),
	}

	filter.AddCommandFlags(result)

	config.AddGitModeFlag(&flags.mode, result.Flags())
	result.Flags().StringVarP(&flags.pattern, "pattern", "p", ".*", "Regexp pattern of the branches to show")

	return result
}
