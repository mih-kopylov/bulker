package branches

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/spf13/cobra"
	"strings"
)

func CreateListCommand() *cobra.Command {
	var filter = runner.Filter{}
	var flags struct {
		mode    gitops.GitMode
		pattern string
	}

	var result = &cobra.Command{
		Use:   "list",
		Short: "Prints a list of repository branches",
		Long: `Prints a list of repository branches.
If a repository doesn't have any branch matching pattern, the repository will be omitted in the result'`,
		RunE: runner.NewDefaultRunner(
			&filter,
			func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				branches, err := gitops.GetBranches(runContext.Fs, runContext.Repo, flags.mode, flags.pattern)
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

	flags.mode = gitops.GitModeAll
	result.Flags().VarP(
		&flags.mode, "mode", "m", fmt.Sprintf(
			"Type of branches to process. "+
				"Available types are: %s, %s, %s", gitops.GitModeAll, gitops.GitModeLocal, gitops.GitModeRemote,
		),
	)
	result.Flags().StringVarP(&flags.pattern, "pattern", "p", ".*", "Regexp pattern of the branches to show")

	return result
}
