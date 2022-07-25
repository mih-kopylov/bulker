package groups

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateCleanCommand() *cobra.Command {
	var result = &cobra.Command{
		Use:   "clean",
		Short: "Removes all configured groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(utils.GetConfiguredFS(), config.ReadConfig())

			sets, err := settingsManager.Read()
			if err != nil {
				return err
			}

			entityInfoMap := map[string]output.EntityInfo{}
			for _, group := range sets.Groups {
				entityInfoMap[group.Name] = output.EntityInfo{Result: "removed", Error: nil}
			}

			sets.Groups = []settings.Group{}

			err = settingsManager.Write(sets)
			if err != nil {
				return err
			}

			err = output.Write("group", entityInfoMap)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return result
}
