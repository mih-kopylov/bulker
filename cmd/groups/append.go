package groups

import (
	"errors"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateAppendCommand() *cobra.Command {
	flags := struct {
		group string
		repos []string
	}{}

	var result = &cobra.Command{
		Use:   "append",
		Short: "Adds repositories to an existing group",
		Long: `Updates the configured group content with adding new repositories. 
If the repo to be added already exists in the group, it will be ignored.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(utils.GetConfiguredFS(), config.ReadConfig())

			sets, err := settingsManager.Read()
			if err != nil {
				return err
			}

			repos, err := utils.GetReposFromStdInOrDefault(flags.repos)
			if err != nil {
				return err
			}

			entityInfoMap := map[string]output.EntityInfo{}

			for _, repoName := range repos {
				err := sets.AddRepoToGroup(flags.group, repoName)
				if err != nil {
					if errors.Is(err, settings.ErrRepoAlreadyAdded) {
						logrus.
							WithField("repo", repoName).
							WithField("group", flags.group).
							Debug("repository already added, skipping")
						entityInfoMap[repoName] = output.EntityInfo{Result: "adding skipped", Error: nil}
					} else {
						entityInfoMap[repoName] = output.EntityInfo{Result: nil, Error: err}
					}
				} else {
					entityInfoMap[repoName] = output.EntityInfo{Result: "added", Error: nil}
				}
			}

			err = settingsManager.Write(sets)
			if err != nil {
				return err
			}

			err = output.Write("repo", entityInfoMap)
			if err != nil {
				return err
			}

			return nil
		},
	}

	result.Flags().StringVarP(&flags.group, "group", "g", "", "Name of the group to update")
	utils.MarkFlagRequiredOrFail(result.Flags(), "group")

	result.Flags().StringSliceVarP(&flags.repos, "repo", "r", []string{}, "Repositories to add to the group")

	utils.AddReadFromStdInFlag(result, "repo")

	return result
}
