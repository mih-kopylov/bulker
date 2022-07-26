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

func CreateUpdateCommand() *cobra.Command {
	flags := struct {
		group         string
		reposToAdd    []string
		reposToRemove []string
	}{}

	var result = &cobra.Command{
		Use:   "update",
		Short: "Updates configured group content",
		Long: `Updates configured group content. 
Runs "add" operation first, and "remove" last.
If the repo to be added already exists in the group, it will be ignored.
If the repo to be removed does not exist in the group, it will be ignored.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(utils.GetConfiguredFS(), config.ReadConfig())

			sets, err := settingsManager.Read()
			if err != nil {
				return err
			}

			entityInfoMap := map[string]output.EntityInfo{}

			for _, repoName := range flags.reposToAdd {
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

			for _, repoName := range flags.reposToRemove {
				err := sets.RemoveRepoFromGroup(flags.group, repoName)
				if err != nil {
					if errors.Is(err, settings.ErrRepoAlreadyRemoved) {
						logrus.
							WithField("repo", repoName).
							WithField("group", flags.group).
							Debug("repository already removed, skipping")
						entityInfoMap[repoName] = output.EntityInfo{Result: "removing skipped", Error: nil}
					} else {
						entityInfoMap[repoName] = output.EntityInfo{Result: nil, Error: err}
					}
				} else {
					entityInfoMap[repoName] = output.EntityInfo{Result: "removed", Error: nil}
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

	result.Flags().StringVarP(&flags.group, "group", "g", "", "Name of the group to remove")
	utils.MarkFlagRequiredOrFail(result.Flags(), "group")

	result.Flags().StringSliceVarP(&flags.reposToAdd, "add", "a", []string{}, "Repositories to add to the group")

	result.Flags().StringSliceVarP(&flags.reposToRemove, "remove", "r", []string{}, "Repositories to remove from the group")

	return result
}
