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

func CreateCheckoutCommand() *cobra.Command {
	var filter = runner.Filter{}
	var flags struct {
		name    string
		discard bool
	}

	var result = &cobra.Command{
		Use:   "checkout",
		Short: "Switch repository to the provided branch",
		RunE: runner.NewDefaultRunner(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Status   gitops.StatusResult
					Checkout gitops.CheckoutResult
					Ref      string
				}

				if flags.discard {
					err := gitops.Discard(runContext.Fs, runContext.Repo)
					if errors.Is(err, fileops.ErrRepositoryNotCloned) {
						return result{gitops.StatusMissing, "", ""}, nil
					}
					if err != nil {
						return nil, fmt.Errorf("failed to discard: %w", err)
					}
				}

				checkoutResult, err := gitops.Checkout(runContext.Fs, runContext.Repo, flags.name)
				if err != nil {
					if errors.Is(err, fileops.ErrRepositoryNotCloned) {
						return result{gitops.StatusMissing, "", ""}, nil
					}
					return nil, fmt.Errorf("failed to checkout: %w", err)
				}

				statusResult, ref, err := gitops.Status(runContext.Fs, runContext.Repo)
				if err != nil {
					return nil, fmt.Errorf("failed to get status: %w", err)
				}

				return result{statusResult, checkoutResult, ref}, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(&flags.name, "branch", "b", "", "Name of the branch to checkout")
	utils.MarkFlagRequiredOrFail(result.Flags(), "branch")

	result.Flags().BoolVarP(
		&flags.discard, "discard", "d", false, "Discards all local changes in the repository before checkout",
	)

	return result
}
