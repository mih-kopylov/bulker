package branches

import (
	"context"
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateCheckoutCommand() *cobra.Command {
	var filter = runner.Filter{}
	var flags struct {
		name string
	}

	var result = &cobra.Command{
		Use:   "checkout",
		Short: "Switch repository to the provided branch",
		RunE: runner.NewDefaultRunner(
			&filter,
			func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Status   string
					Checkout string
					Ref      string
				}

				checkoutStatus, err := gitops.Checkout(runContext.Fs, runContext.Repo, flags.name)
				if err != nil {
					if errors.Is(err, gitops.ErrRepositoryNotCloned) {
						return result{"MISSING", "", ""}, nil
					}
					return nil, fmt.Errorf("failed to checkout: %w", err)
				}

				repoStatus, ref, err := gitops.Status(runContext.Fs, runContext.Repo)
				if err != nil {
					return nil, fmt.Errorf("failed to get status: %w", err)
				}

				switch repoStatus {
				case gitops.StatusOk:
					return result{"OK", checkoutStatus.String(), ref}, nil
				case gitops.StatusDirty:
					return result{"DIRTY", checkoutStatus.String(), ref}, nil
				case gitops.StatusMissing:
					return result{"MISSING", checkoutStatus.String(), ref}, nil
				default:
					return nil, fmt.Errorf("unsupported status %v", repoStatus)
				}
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(&flags.name, "branch", "b", "", "Name of the branch to checkout")
	utils.MarkFlagRequiredOrFail(result.Flags(), "branch")

	return result
}
