package groups

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateGetCommand() *cobra.Command {
	flags := struct {
		group string
	}{}

	var result = &cobra.Command{
		Use:   "get",
		Short: "Prints repositories of the provided group",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(config.ReadConfig())

			sets, err := settingsManager.Read()
			if err != nil {
				return err
			}

			group, err := sets.GetGroup(flags.group)
			if err != nil {
				return err
			}

			entityInfoMap := map[string]output.EntityInfo{}
			for _, repoName := range group.Repos {
				entityInfoMap[repoName] = output.EntityInfo{Result: nil, Error: nil}
			}
			err = output.Write(cmd.OutOrStdout(), "repo", entityInfoMap)
			if err != nil {
				return err
			}

			return nil
		},
	}

	result.Flags().StringVarP(&flags.group, "group", "g", "", "Name of the group to get")
	utils.MarkFlagRequiredOrFail(result.Flags(), "group")

	return result
}
