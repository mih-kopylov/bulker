package files

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

func CreateSearchCommand() *cobra.Command {
	var filter = runner.Filter{}
	var flags = struct {
		pattern     string
		contains    string
		replacement string
		before      int
		after       int
	}{}

	var result = &cobra.Command{
		Use:   "search",
		Short: "Searches for files and file content",
		RunE: runner.NewDefaultRunner(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				searchResult, err := fileops.SearchFiles(
					runContext.Fs, runContext.Repo, flags.pattern, flags.contains,
					flags.before, flags.after,
				)
				if err != nil {
					if errors.Is(err, fileops.ErrSourceNotFound) {
						return nil, nil
					}
					return nil, err
				}

				if len(searchResult) == 0 {
					return nil, nil
				}

				if flags.contains == "" {
					// result without files content, just found files names
					buffer := bytes.Buffer{}
					for _, item := range searchResult {
						relFileName, err := filepath.Rel(runContext.Repo.Path, item.FileName)
						if err != nil {
							return nil, err
						}

						buffer.WriteString(fmt.Sprintf("%s\n", relFileName))
					}

					return buffer.String(), nil
				} else {
					// result with files content
					result := make(map[string][]string)
					for _, item := range searchResult {
						relFileName, err := filepath.Rel(runContext.Repo.Path, item.FileName)
						if err != nil {
							return nil, err
						}

						result[relFileName] = item.Matches
					}

					resultBytes, err := yaml.Marshal(result)
					if err != nil {
						return nil, err
					}

					return string(resultBytes), nil
				}
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(
		&flags.pattern, "files", "f", "*",
		`Glob files pattern to search. If not passed, all files will be processed.
See https://pkg.go.dev/path/filepath#Match for syntax`,
	)

	result.Flags().StringVarP(
		&flags.contains, "contains", "c", "",
		`Regexp to search in the files for. 
If passed, only repos with matching files will be returned.
Example: "text\s+(value)"`,
	)

	result.Flags().StringVarP(
		&flags.replacement, "replacement", "r", "",
		`Expression to use for replacement.
If passed, the matched files content will be updated.
Only matters if "--contains" is set.
Example: "$1 found"`,
	)

	result.Flags().IntVarP(
		&flags.before, "before", "b", 0,
		`Number of preceding rows to print in addition to the matching ones. 
Only matters if "--pattern" and "--contains" are set.`,
	)

	result.Flags().IntVarP(
		&flags.after, "after", "a", 0,
		`Number of following rows to print in addition to the matching ones. 
Only matters if "--pattern" and "--contains" are set.`,
	)

	return result
}
