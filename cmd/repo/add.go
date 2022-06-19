package repo

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new repository to the supported list",
	RunE: func(cmd *cobra.Command, args []string) error {
		settingsManager := settings.NewManager(afero.NewOsFs())

		sets, err := settingsManager.Read()
		if err != nil {
			return err
		}

		err = sets.AddRepo(addFlags.name, addFlags.url, addFlags.tags)
		if err != nil {
			return err
		}

		err = settingsManager.Write(sets)
		if err != nil {
			return err
		}

		logrus.WithField("repo", addFlags.name).Info("repository added")

		return nil
	},
}

var addFlags struct {
	name string
	url  string
	tags []string
}

func init() {
	AddCmd.Flags().StringVar(&addFlags.name, "name", "", "Name of the repository")
	utils.MarkFlagRequiredOrFail(AddCmd.Flags(), "name")

	AddCmd.Flags().StringVar(&addFlags.url, "url", "", "URL of the repository")
	utils.MarkFlagRequiredOrFail(AddCmd.Flags(), "url")

	AddCmd.Flags().StringSliceVar(&addFlags.tags, "tags", []string{}, "Tags of the repository")
}
