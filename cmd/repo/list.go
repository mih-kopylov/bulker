package repo

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints a list of supported repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		newRunner := runner.NewRunner(utils.GetConfiguredFS(), config.ReadConfig(), listFilter)

		err := newRunner.Run(
			func(ctx context.Context, runContext *runner.RunContext) (runner.Result, error) {
				return listResult{Message: runContext.Repo.Name}, nil
			},
		)
		if err != nil {
			return err
		}

		return nil
	},
}

type listResult struct {
	Message string
}

var listFilter = &runner.Filter{}

func init() {
	listFilter.AddCommandFlags(ListCmd)
}
