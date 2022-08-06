package git

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/spf13/cobra"
)

func CreateFetchCommand() *cobra.Command {
	var filter = runner.Filter{}

	var result = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch changes from remote",
		RunE: runner.NewDefaultRunner(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				err := gitops.Fetch(runContext.Fs, runContext.Repo)
				if err != nil {
					return nil, err
				}

				return "fetched successfully", nil
			},
		),
	}

	filter.AddCommandFlags(result)

	return result
}
