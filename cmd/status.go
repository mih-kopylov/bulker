package cmd

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateStatusCommand() *cobra.Command {
	var filter = runner.Filter{}

	flags := struct {
		showOk      bool
		showDirty   bool
		showMissing bool
	}{}

	var result = &cobra.Command{
		Use:   "status",
		Short: "Prints status of all registered repositories",
		Long: `Prints status of all registered repositories. Status value can be one of:
* OK - the repository successfully cloned, there are no uncommitted changes
* DIRTY - the repository successfully cloned, but there are uncommitted changes
* MISSING - the repository is not cloned yet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			newRunner, err := runner.NewRunner(utils.GetConfiguredFS(), config.ReadConfig(), &filter)
			if err != nil {
				return err
			}

			err = newRunner.Run(
				func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
					repoStatus, err := gitops.Status(runContext.FS, runContext.Repo)
					if err != nil {
						return nil, fmt.Errorf("failed to get status: %w", err)
					}

					switch repoStatus {
					case gitops.StatusOk:
						if flags.showOk || (!flags.showOk && !flags.showDirty && !flags.showMissing) {
							return "OK", nil
						}
						return nil, nil
					case gitops.StatusDirty:
						if flags.showDirty || (!flags.showOk && !flags.showDirty && !flags.showMissing) {
							return "DIRTY", nil
						}
						return nil, nil
					case gitops.StatusMissing:
						if flags.showMissing || (!flags.showOk && !flags.showDirty && !flags.showMissing) {
							return "MISSING", nil
						}
						return nil, nil
					default:
						return nil, fmt.Errorf("unsupported status %v", repoStatus)
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

	result.Flags().BoolVar(&flags.showOk, "ok", false, "Keep repositories with 'OK' status")
	result.Flags().BoolVar(&flags.showDirty, "dirty", false, "Keep repositories with 'DIRTY' status")
	result.Flags().BoolVar(&flags.showMissing, "missing", false, "Keep repositories with 'MISSING' status")

	return result
}
