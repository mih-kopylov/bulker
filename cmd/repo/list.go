package repo

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints a list of supported repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		settingsManager := settings.NewManager(afero.NewOsFs())

		sets, err := settingsManager.Read()
		if err != nil {
			return err
		}

		for _, repo := range sets.Repos {
			logrus.Info(repo.Name)
		}

		return nil
	},
}
