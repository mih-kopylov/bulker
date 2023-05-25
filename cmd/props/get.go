package props

import (
	"context"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"strings"

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
		Long: `Find property value in matching files.

If a repository doesn't contain any file matching file pattern, the repository won't be printed.
If a repository contains files matching file pattern, but none of the files contain path pattern, 
the repository won't be printed.

Examples.

The command below will print verison of the parent pom in all repositories that have pom.xml in root directory:

    bulker props get -f pom.xml -p $.parent.version

The command below will search for Spring applications and print their names: 

    bulker props get -f src/**/application.yaml -p $.spring.application.name

`,
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Result string
				}

				fileSearchResults, err := fileops.SearchFiles(runContext.Repo, flags.filePattern, "", 0, 0)
				if err != nil {
					return nil, fmt.Errorf("failed to search files with properties: %w", err)
				}

				logrus.WithField("repo", runContext.Repo.Name).
					WithField("count", len(fileSearchResults)).
					Debug("found files")

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

				if len(foundFiles) == 0 {
					return nil, nil
				}

				marshalled, err := yaml.Marshal(foundFiles)
				if err != nil {
					return nil, errors.Wrap(err, "failed to marshall found files")
				}
				foundFilesString := strings.TrimSpace(string(marshalled))

				return result{Result: foundFilesString}, nil
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
	utils.MarkFlagRequiredOrFail(result.Flags(), "path")

	return result
}
