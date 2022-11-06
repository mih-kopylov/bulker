package branches

import (
	"context"
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateCreateCommand() *cobra.Command {
	var filter = runner.Filter{}
	var flags struct {
		name    string
		discard bool
	}

	var result = &cobra.Command{
		Use:   "create",
		Short: "Create a new branch and switch repository to it",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Status gitops.StatusResult
					Create gitops.CreateResult
					Ref    string
				}

				if flags.discard {
					err := gitops.Discard(runContext.Repo)
					if errors.Is(err, fileops.ErrRepositoryNotCloned) {
						return result{gitops.StatusMissing, "", ""}, nil
					}
					if err != nil {
						return nil, fmt.Errorf("failed to discard: %w", err)
					}
				}

				createResult, err := gitops.CreateBranch(runContext.Repo, flags.name)
				if err != nil {
					if errors.Is(err, fileops.ErrRepositoryNotCloned) {
						return result{gitops.StatusMissing, "", ""}, nil
					}
					return nil, fmt.Errorf("failed to create: %w", err)
				}

				_, err = gitops.Checkout(runContext.Repo, flags.name)
				if err != nil {
					if errors.Is(err, fileops.ErrRepositoryNotCloned) {
						return result{gitops.StatusMissing, "", ""}, nil
					}
					return nil, fmt.Errorf("failed to checkout: %w", err)
				}

				statusResult, ref, err := gitops.Status(runContext.Repo)
				if err != nil {
					return nil, fmt.Errorf("failed to get status: %w", err)
				}

				return result{statusResult, createResult, ref}, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(&flags.name, "branch", "b", "", "Name of the branch to create")
	utils.MarkFlagRequiredOrFail(result.Flags(), "branch")

	result.Flags().BoolVarP(
		&flags.discard, "discard", "d", false, "Discards all local changes in the repository before creating",
	)

	return result
}
