package cmd

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateConfigureCommand() *cobra.Command {
	flags := struct {
		configFileName string
		gitMode        config.GitMode
	}{}

	var result = &cobra.Command{
		Use:   "configure",
		Short: "Configures bulker and saves the configuration to file for future calls",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := config.ReadConfig()
			conf.GitMode = flags.gitMode
			err := config.WriteConfig(conf, flags.configFileName)
			if err != nil {
				return err
			}

			err = output.Write(
				cmd.OutOrStdout(), "entity",
				map[string]output.EntityInfo{
					"configuration": {Result: "saved"},
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	result.Flags().StringVarP(&flags.configFileName, "save", "s", utils.AbsPathify("$HOME/.bulker/bulker.yaml"), "")
	config.AddGitModeFlag(&flags.gitMode, result.Flags())

	return result
}
