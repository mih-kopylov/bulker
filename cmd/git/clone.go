package git

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/spf13/cobra"
)

func CreateCloneCommand() *cobra.Command {
	var filter = runner.Filter{}
	var flags = struct {
		recreate bool
	}{}

	var result = &cobra.Command{
		Use:   "clone",
		Short: "Clones the configured repositories out if they have not been yet",
		RunE: runner.NewCommandRunner(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				cloneResult, err := gitops.CloneRepo(runContext.Repo, flags.recreate)
				if err != nil {
					return nil, fmt.Errorf("failed to clone: %w", err)
				}

				return cloneResult.String(), nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().BoolVar(
		&flags.recreate, "recreate", false,
		"Delete the repository directory and start clone from scratch",
	)

	return result
}
