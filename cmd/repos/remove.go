package repos

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateRemoveCommand(sh shell.Shell) *cobra.Command {
	var flags struct {
		name string
	}

	var result = &cobra.Command{
		Use:   "remove",
		Short: "Remove one repo from the supported list",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(config.ReadConfig(), sh)

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
				cmd.OutOrStdout(), "repo",
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

	result.Flags().StringVarP(&flags.name, "name", "n", "", "Name of the repository")
	utils.MarkFlagRequiredOrFail(result.Flags(), "name")

	return result
}
