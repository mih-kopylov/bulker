package git

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreatePullCommand() *cobra.Command {
	var filter = runner.Filter{}

	var result = &cobra.Command{
		Use:   "pull",
		Short: "Pull changes from remote origin",
		RunE: func(cmd *cobra.Command, args []string) error {
			newRunner, err := runner.NewRunner(utils.GetConfiguredFS(), config.ReadConfig(), &filter)
			if err != nil {
				return err
			}

			err = newRunner.Run(
				func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
					err := gitops.Pull(runContext.FS, runContext.Repo)
					if err != nil {
						return nil, err
					}

					return "pulled successfully", nil
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
