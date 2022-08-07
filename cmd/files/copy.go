package files

import (
	"context"
	"errors"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateCopyCommand() *cobra.Command {
	var filter = runner.Filter{}
	var flags = struct {
		source string
		target string
	}{}

	var result = &cobra.Command{
		Use:   "copy",
		Short: "Copy source file into target",
		RunE: runner.NewDefaultRunner(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Status string
					Source string
					Target string
				}

				source, target, err := fileops.Copy(runContext.Fs, runContext.Repo, flags.source, flags.target)
				if err != nil {
					if errors.Is(err, fileops.ErrSourceNotFound) {
						return nil, nil
					}
					return result{"failed", source, target}, err
				}

				return result{"copied", source, target}, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVar(
		&flags.source, "source", "",
		`File name to copy. The path is considered relative to the repository root.
The following examples are equal:
- README.md
- ./README.md
- /README.md`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "source")

	result.Flags().StringVar(
		&flags.target, "target", "",
		`File name to copy to. The same rules as for "source" flag are applied`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "target")

	return result
}
