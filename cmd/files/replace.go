package files

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
	"path/filepath"
)

func CreateReplaceCommand() *cobra.Command {
	var filter = runner.Filter{}
	var flags = struct {
		pattern     string
		contains    string
		replacement string
	}{}

	var result = &cobra.Command{
		Use:   "replace",
		Short: "Replaces files content",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				replaceResult, err := fileops.ReplaceInFiles(
					runContext.Repo, flags.pattern, flags.contains, flags.replacement,
				)
				if err != nil {
					if errors.Is(err, fileops.ErrSourceNotFound) {
						return nil, nil
					}
					return nil, err
				}

				if len(replaceResult) == 0 {
					return nil, nil
				}

				buffer := bytes.Buffer{}
				for _, replacementResult := range replaceResult {
					relFileName, err := filepath.Rel(runContext.Repo.Path, replacementResult.FileName)
					if err != nil {
						return nil, err
					}

					buffer.WriteString(fmt.Sprintf("%v :: %v\n", relFileName, replacementResult.Count))
				}

				return buffer.String(), nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(
		&flags.pattern, "files", "f", "*",
		`Glob files pattern to search. If not passed, all files will be processed.
See https://pkg.go.dev/path/filepath#Match for syntax`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "files")

	result.Flags().StringVarP(
		&flags.contains, "contains", "c", "",
		`Regexp to search in the files for. 
Example: "text\s+(value)"`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "contains")

	result.Flags().StringVarP(
		&flags.replacement, "replacement", "r", "",
		`Value to insert instead of findings`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "replacement")

	return result
}
