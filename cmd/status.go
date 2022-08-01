package cmd

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
	"strings"
)

func CreateStatusCommand() *cobra.Command {
	var filter = runner.Filter{}

	flags := struct {
		status string
		ref    string
	}{}

	var result = &cobra.Command{
		Use:   "status",
		Short: "Prints status of all registered repositories",
		Long: `Prints status of all registered repositories. Status value can be one of:
* Clean - the repository successfully cloned, there are no uncommitted changes
* Dirty - the repository successfully cloned, but there are uncommitted changes
* Missing - the repository is not cloned yet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			newRunner, err := runner.NewRunner(utils.GetConfiguredFS(), config.ReadConfig(), &filter)
			if err != nil {
				return err
			}

			err = newRunner.Run(
				func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
					repoStatus, ref, err := gitops.Status(runContext.Fs, runContext.Repo)
					if err != nil {
						return nil, fmt.Errorf("failed to get status: %w", err)
					}

					type result struct {
						Status string
						Ref    string
					}

					statusFlagMatches := func(repoStatus gitops.StatusResult) bool {
						if flags.status == "" {
							return true
						}

						hasPrefix := strings.HasPrefix(flags.status, "!")

						if hasPrefix && !strings.EqualFold(flags.status[1:], repoStatus.String()) {
							return true
						}

						if !hasPrefix && strings.EqualFold(flags.status, repoStatus.String()) {
							return true
						}

						return false
					}

					refFlagMatches := func(repoRef string) bool {
						if flags.ref == "" {
							return true
						}

						hasPrefix := strings.HasPrefix(flags.ref, "!")

						if hasPrefix && !strings.EqualFold(flags.ref[1:], repoRef) {
							return true
						}

						if !hasPrefix && strings.EqualFold(flags.ref, repoRef) {
							return true
						}

						return false
					}

					if statusFlagMatches(repoStatus) && refFlagMatches(ref) {
						return result{repoStatus.String(), ref}, nil
					}
					return nil, nil
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVar(
		&flags.status, "status", "", fmt.Sprintf(
			`Keep repositories with specified status.
Examples: 
"bulker status --status clean" - will keep only repositories with "%v" status
"bulker status --status !clean" - will keep repositories with any status except "%v"`, gitops.StatusClean,
			gitops.StatusClean,
		),
	)
	result.Flags().StringVar(
		&flags.ref, "ref", "",
		`Keep repositories with specified ref.
Examples: 
"bulker status --ref master" - will keep only repositories with "master" ref
"bulker status --ref !master" - will keep repositories with any ref except "master"`,
	)

	return result
}
