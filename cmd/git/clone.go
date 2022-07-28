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
	var filter = runner.Filter{}
	var flags = struct {
		recreate bool
	}{}

	var result = &cobra.Command{
		Use:   "clone",
		Short: "Clones the configured repositories out if they have not been yet",
		RunE: func(cmd *cobra.Command, args []string) error {
			newRunner, err := runner.NewRunner(utils.GetConfiguredFS(), config.ReadConfig(), &filter)
			if err != nil {
				return err
			}

			err = newRunner.Run(
				func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
					cloneResult, err := gitops.CloneRepo(runContext.FS, runContext.Repo, flags.recreate)
					if err != nil {
						return nil, fmt.Errorf("failed to clone: %w", err)
					}

					switch cloneResult {
					case gitops.ClonedSuccessfully:
						return "cloned successfully", nil
					case gitops.ClonedAgain:
						return "cloned again", nil
					case gitops.ClonedAlready:
						return "already cloned", nil
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

	result.Flags().BoolVar(
		&flags.recreate, "recreate", false,
		"Delete the repository directory and start clone from scratch",
	)

	return result
}
