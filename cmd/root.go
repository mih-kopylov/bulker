package cmd

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateRootCommand() *cobra.Command {
	var result = &cobra.Command{
		Use:   "bulker",
		Short: "Runs different operations on a bunch of repositories in bulk mode",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configureLogrus()
		},
	}

	result.PersistentFlags().Bool("debug", false, "Enable debug level logging")
	utils.BindFlag(result.PersistentFlags().Lookup("debug"), "debug")

	result.PersistentFlags().String(
		"settings", utils.AbsPathify("$HOME/.bulker/settings.yaml"),
		"Settings file name, where list of repositories is stored",
	)
	utils.BindFlag(result.PersistentFlags().Lookup("settings"), "settings")

	result.PersistentFlags().String(
		"output", string(config.LogOutputFormat), fmt.Sprintf(
			"Set commands output format. Available formats: %v, %v, %v", config.LogOutputFormat,
			config.LineOutputFormat, config.JsonOutputFormat,
		),
	)
	utils.BindFlag(result.PersistentFlags().Lookup("output"), "output")

	result.AddCommand(CreateReposCommand())
	result.AddCommand(CreateGitCommand())
	result.AddCommand(CreateGroupsCommand())

	return result
}

func Execute() {
	rootCmd := CreateRootCommand()
	err := rootCmd.Execute()
	if err != nil {
		logrus.Debugf("command failed: %v", err)
	}
}

func init() {
	configureViper()
}

func configureViper() {
	viper.SetConfigName("bulker")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(utils.AbsPathify("$HOME/.bulker"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// ignore case when config file is not found
		} else {
			panic(fmt.Errorf("can't read config: %w", err))
		}
	}
	logrus.WithField("file", viper.ConfigFileUsed()).Debug("config used")
}

func configureLogrus() {
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugf("debug logging enabled")
	}
}
