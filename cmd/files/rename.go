package files

import (
	"context"
	"errors"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateRenameCommand() *cobra.Command {
	var filter = runner.Filter{}
	var flags = struct {
		source string
		target string
		force  bool
	}{}

	var result = &cobra.Command{
		Use:   "rename",
		Short: "Rename source file into target",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Status string
					Source string
					Target string
				}

				source, target, err := fileops.Rename(
					runContext.Fs, runContext.Repo, flags.source, flags.target,
					flags.force,
				)
				if err != nil {
					if errors.Is(err, fileops.ErrSourceNotFound) {
						return nil, nil
					}
					return result{"failed", source, target}, err
				}

				return result{"renamed", source, target}, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVar(
		&flags.source, "source", "",
		`File name to rename. The path is considered relative to the repository root.
The following examples are equal:
- README.md
- ./README.md
- /README.md`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "source")

	result.Flags().StringVar(
		&flags.target, "target", "",
		`File name to rename to. The same rules as for "source" flag are applied`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "target")

	result.Flags().BoolVarP(
		&flags.force, "force", "f", false,
		"Force renaming if target file already exists, rewriting the target file",
	)

	return result
}
