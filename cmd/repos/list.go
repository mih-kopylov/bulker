package repos

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateListCommand() *cobra.Command {
	var filter = &runner.Filter{}

	var result = &cobra.Command{
		Use:   "list",
		Short: "Prints a list of supported repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			newRunner := runner.NewRunner(utils.GetConfiguredFS(), config.ReadConfig(), filter)

			err := newRunner.Run(
				func(ctx context.Context, runContext *runner.RunContext) (runner.Result, error) {
					return runContext.Repo.Name, nil
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	filter.AddCommandFlags(result)

	return result
}
