package files

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
	"strings"
)

func CreateRemoveCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}
	var flags = struct {
		pattern string
	}{}

	var result = &cobra.Command{
		Use:   "remove",
		Short: "Remove all files matching pattern",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Removed string
				}

				removed, err := fileops.Remove(runContext.Repo, flags.pattern)
				if len(removed) == 0 {
					return nil, nil
				}

				removedString := strings.Join(removed, "\n")
				if err != nil {
					return result{removedString}, err
				}

				return result{removedString}, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(
		&flags.pattern, "files", "f", "",
		`Glob files pattern to remove.
See https://pkg.go.dev/path/filepath#Match for syntax`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "files")

	return result
}
