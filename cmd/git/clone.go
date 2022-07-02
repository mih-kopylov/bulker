package git

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateCloneCommand() *cobra.Command {
	type CloneResult struct {
		Message string
	}
	var filter = &runner.Filter{}

	var result = &cobra.Command{
		Use:   "clone",
		Short: "Clones the configured repositories out if they have not been yet",
		RunE: func(cmd *cobra.Command, args []string) error {
			newRunner := runner.NewRunner(utils.GetConfiguredFS(), config.ReadConfig(), filter)

			err := newRunner.Run(
				func(ctx context.Context, runContext *runner.RunContext) (runner.Result, error) {
					cloneResult, err := gitops.CloneRepo(runContext.FS, runContext.Repo)
					if err != nil {
						return nil, err
					}

					switch cloneResult {
					case gitops.ClonedSuccessfully:
						return &CloneResult{"cloned successfully"}, nil
					case gitops.AlreadyCloned:
						return &CloneResult{"already cloned"}, nil
					default:
						return nil, fmt.Errorf("unsupported clone status: status=%v", cloneResult)
					}
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
