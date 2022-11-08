package groups

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/spf13/cobra"
)

func CreateListCommand() *cobra.Command {
	var result = &cobra.Command{
		Use:   "list",
		Short: "Prints a list of configured groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(config.ReadConfig())

			sets, err := settingsManager.Read()
			if err != nil {
				return err
			}

			entityInfoMap := map[string]output.EntityInfo{}
			for _, group := range sets.Groups {
				entityInfoMap[group.Name] = output.EntityInfo{Result: nil, Error: nil}
			}

			err = output.Write(cmd.OutOrStdout(), "group", entityInfoMap)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return result
}
