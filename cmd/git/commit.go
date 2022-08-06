package git

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateCommitCommand() *cobra.Command {
	var filter = runner.Filter{}

	var flags struct {
		message string
		pattern string
	}

	var result = &cobra.Command{
		Use:   "commit",
		Short: "Commit changes",
		RunE: runner.NewDefaultRunner(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				err := gitops.Commit(runContext.Fs, runContext.Repo, flags.pattern, flags.message)
				if err != nil {
					return nil, err
				}

				return "committed", nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(&flags.message, "message", "m", "", "Commit message")
	utils.MarkFlagRequiredOrFail(result.Flags(), "message")

	result.Flags().StringVarP(
		&flags.pattern, "pattern", "p", "**", `File pattern to commit.
See https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-glob for documentation.
If missed, all added/changed/removed files will be committed.'`,
	)

	return result
}
