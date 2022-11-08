package groups

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
		group string
	}

	var result = &cobra.Command{
		Use:   "remove",
		Short: "Removes a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(config.ReadConfig(), sh)

			sets, err := settingsManager.Read()
			if err != nil {
				return err
			}

			err = sets.RemoveGroup(flags.group)
			if err != nil {
				return err
			}

			err = settingsManager.Write(sets)
			if err != nil {
				return err
			}

			entityInfoMap := map[string]output.EntityInfo{}
			entityInfoMap[flags.group] = output.EntityInfo{Result: "removed", Error: nil}
			err = output.Write(cmd.OutOrStdout(), "group", entityInfoMap)
			if err != nil {
				return err
			}

			return nil
		},
	}

	result.Flags().StringVarP(&flags.group, "group", "g", "", "Name of the group to remove")
	utils.MarkFlagRequiredOrFail(result.Flags(), "group")

	return result
}
