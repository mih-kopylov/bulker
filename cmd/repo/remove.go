package repo

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove one repo from the supported list",
	RunE: func(cmd *cobra.Command, args []string) error {
		settingsManager := settings.NewManager(utils.GetConfiguredFS(), config.ReadConfig())

		sets, err := settingsManager.Read()
		if err != nil {
			return err
		}

		err = sets.RemoveRepo(removeFlags.name)
		if err != nil {
			return err
		}

		err = settingsManager.Write(sets)
		if err != nil {
			return err
		}

		logrus.WithField("repo", removeFlags.name).Info("repository removed")

		return nil
	},
}

var removeFlags struct {
	name string
}

func init() {
	RemoveCmd.Flags().StringVar(&removeFlags.name, "name", "", "Name of the repository")
	utils.MarkFlagRequiredOrFail(AddCmd.Flags(), "name")
}
