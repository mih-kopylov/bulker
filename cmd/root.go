package cmd

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/signal"
	"syscall"
)

func CreateRootCommand(applicationVersion string, sh shell.Shell) *cobra.Command {
	var result = &cobra.Command{
		Use:     "bulker",
		Short:   "Runs different operations on a bunch of repositories in bulk mode",
		Version: applicationVersion,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configureLogrus()
		},
	}

	result.SetVersionTemplate("{{.Version}}")

	// once the APP gets a signal, it will mark the context as Done
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	result.SetContext(ctx)

	result.PersistentFlags().Bool("debug", false, "Enable debug level logging. Hide progress bar as well.")
	utils.BindFlag(result.PersistentFlags().Lookup("debug"), "debug")

	result.PersistentFlags().String(
		"settings", utils.AbsPathify("$HOME/.bulker/settings.yaml"),
		"Settings file name, where list of repositories and other user data is stored",
	)
	utils.BindFlag(result.PersistentFlags().Lookup("settings"), "settingsFileName")

	result.PersistentFlags().String(
		"repos-directory", utils.AbsPathify("."),
		"Root directory for all the repositories to store",
	)
	utils.BindFlag(result.PersistentFlags().Lookup("repos-directory"), "reposDirectory")

	var runMode = config.Parallel
	result.PersistentFlags().Var(
		&runMode,
		"run-mode",
		"Parallel (par) or sequential (seq) run mode for repositories processing",
	)
	utils.BindFlag(result.PersistentFlags().Lookup("run-mode"), "runMode")

	result.PersistentFlags().String(
		"max-workers", "10",
		"Maximum number of workers to process repositories simultaneously",
	)
	utils.BindFlag(result.PersistentFlags().Lookup("max-workers"), "maxWorkers")

	result.PersistentFlags().Bool(
		"no-progress", false,
		"Do not show progress bar during repositories processing",
	)
	utils.BindFlag(result.PersistentFlags().Lookup("no-progress"), "noProgress")

	result.PersistentFlags().String(
		"output", string(config.TableOutputFormat), fmt.Sprintf(
			"Set commands output format. Available formats: %v, %v, %v, %v", config.LogOutputFormat,
			config.LineOutputFormat, config.JsonOutputFormat, config.TableOutputFormat,
		),
	)
	utils.BindFlag(result.PersistentFlags().Lookup("output"), "output")

	result.AddCommand(CreateReposCommand(sh))
	result.AddCommand(CreateGitCommand(sh))
	result.AddCommand(CreateGroupsCommand(sh))
	result.AddCommand(CreateStatusCommand(sh))
	result.AddCommand(CreateRunCommand(sh))
	result.AddCommand(CreateOpenCommand(sh))
	result.AddCommand(CreateFilesCommand(sh))
	result.AddCommand(CreateConfigureCommand())

	return result
}

func init() {
	configureViper()
}

func configureViper() {
	viper.SetConfigName("bulker")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(utils.AbsPathify("$HOME/.bulker"))
	viper.SetEnvPrefix("B")
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
