package props

import (
	"context"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"fmt"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/mih-kopylov/bulker/pkg/props"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"path/filepath"
)

func CreateGetCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}
	var flags = struct {
		filePattern string
		path        string
	}{}

	var result = &cobra.Command{
		Use:   "get",
		Short: "Find property value in matching files",
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Files string
				}

				fileSearchResults, err := fileops.SearchFiles(runContext.Repo, flags.filePattern, "", 0, 0)
				if err != nil {
					return nil, fmt.Errorf("failed to search files with properties: %w", err)
				}

				logrus.WithField("repo", runContext.Repo.Name).
					WithField("count", len(fileSearchResults)).
					Debug("found files")

				if len(fileSearchResults) == 0 {
					return nil, nil
				}

				foundFiles := map[string]string{}

				for _, searchResult := range fileSearchResults {
					shortFileName, err := filepath.Rel(runContext.Repo.Path, searchResult.FileName)
					if err != nil {
						return nil, fmt.Errorf("failed to get relative path of %s: %w", searchResult.FileName, err)
					}

					property, err := props.GetPropertyFromFile(searchResult.FileName, flags.path)
					if err != nil {
						if errors.Is(err, props.ErrPropertyNotFound) {
							logrus.WithField("repo", runContext.Repo.Path).
								WithField("file", shortFileName).
								Debug("property not found")
							continue
						}
						return nil, err
					}

					foundFiles[shortFileName] = property
				}

				marshalled, err := yaml.Marshal(foundFiles)
				if err != nil {
					return nil, errors.Wrap(err, "failed to marshall found files")
				}

				return result{Files: string(marshalled)}, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().StringVarP(
		&flags.filePattern, "files", "f", "",
		`Glob files pattern to search properties for.
See https://pkg.go.dev/path/filepath#Match for syntax`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "files")

	result.Flags().StringVarP(
		&flags.path, "path", "p", "",
		`Property path. Similar to JsonPath, but very simplified. Dot separated chain of nodes, 
starts with "$" as a root.
Example: "$.app.reload.enabled"`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "files")

	return result
}
