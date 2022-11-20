package files

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
)

func CreateSearchCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}
	var flags = struct {
		pattern  string
		contains string
		before   int
		after    int
	}{}

	var result = &cobra.Command{
		Use:   "search",
		Short: "Searches for files and file content",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				searchResult, err := fileops.SearchFiles(
					runContext.Repo, flags.pattern, flags.contains, flags.before, flags.after,
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

					return strings.TrimSpace(buffer.String()), nil
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

					w := &bytes.Buffer{}
					encoder := yaml.NewEncoder(w)
					encoder.SetIndent(2)
					err := encoder.Encode(result)
					if err != nil {
						return nil, err
					}

					resultBytes, err := yaml.Marshal(result)
					if err != nil {
						return nil, err
					}

					return strings.TrimSpace(string(resultBytes)), nil
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
	utils.MarkFlagRequiredOrFail(result.Flags(), "files")

	result.Flags().StringVarP(
		&flags.contains, "contains", "c", "",
		`Regexp to search in the files for. 
If passed, only repos with matching files will be returned.
Example: "text\s+(value)"`,
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
