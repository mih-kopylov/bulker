package repos

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateRemoveCommand() *cobra.Command {
	var flags struct {
		name string
	}

	var result = &cobra.Command{
		Use:   "remove",
		Short: "Remove one repo from the supported list",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(utils.GetConfiguredFS(), config.ReadConfig())

			sets, err := settingsManager.Read()
			if err != nil {
				return err
			}

			err = sets.RemoveRepo(flags.name)
			if err != nil {
				return err
			}

			err = settingsManager.Write(sets)
			if err != nil {
				return err
			}

			err = output.Write(
				"repo",
				map[string]output.EntityInfo{
					flags.name: {Result: "removed"},
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	result.Flags().StringVar(&flags.name, "name", "", "Name of the repository")
	utils.MarkFlagRequiredOrFail(result.Flags(), "name")

	return result
}
