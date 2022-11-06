package cmd

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateOpenCommand() *cobra.Command {
	var filter = runner.Filter{}

	var flags struct {
		page shell.RepoPage
	}

	var result = &cobra.Command{
		Use:   "open",
		Short: "Opens repository page in a browser",
		Long: fmt.Sprintf(
			`Opens repository page in a browser.
It's up to the OS to choose which browser is used to open the URL.

The following VCS platforms are supported:
- %v
- %v`, shell.RepoTypeNameGithubCom, shell.RepoTypeNameBitbucketOrg,
		),
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Status string
					Url    string
				}

				url, err := shell.OpenPage(runContext.Repo, flags.page)
				if err != nil {
					return nil, fmt.Errorf("failed to open: %w", err)
				}

				return result{"opened", url}, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	flags.page = shell.PageSources
	result.Flags().VarP(
		&flags.page, "page", "p", fmt.Sprintf(
			`A page to open. 
The following pages are supported: %s %s %s %s`, shell.PageSources.Name, shell.PageBranches.Name,
			shell.PageBuilds.Name, shell.PagePulls.Name,
		),
	)

	return result
}
