package repos

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
	"strings"
)

func CreateListCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}

	var result = &cobra.Command{
		Use:   "list",
		Short: "Prints a list of supported repositories",
		RunE: runner.NewCommandRunner(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				type result struct {
					Url  string
					Tags string
				}

				settingsManager := settings.NewManager(config.ReadConfig(), sh)
				sets, err := settingsManager.Read()
				if err != nil {
					return nil, err
				}

				repo, err := sets.GetRepo(runContext.Repo.Name)
				if err != nil {
					return nil, err
				}

				return result{Url: repo.Url, Tags: strings.Join(repo.Tags, ", ")}, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	return result
}
