package cmd

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "bulker",
	Short: "Runs different operations on a bunch of repositories in bulk mode",
}

func Execute() {
	rootCmd.AddCommand(repoCmd)
	rootCmd.AddCommand(gitCmd)

	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatalf("command failed: %v", err)
	}
}

func init() {
	configureViper()

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug level logging")
	utils.BindFlag(rootCmd.PersistentFlags().Lookup("debug"), "debug")

	rootCmd.PersistentFlags().String(
		"settings", utils.AbsPathify("$HOME/.bulker/settings.yaml"),
		"Settings file name, where list of repositories is stored",
	)
	utils.BindFlag(rootCmd.PersistentFlags().Lookup("settings"), "settings")

	configureLogrus()
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
			//ignore case when config file is not found
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
