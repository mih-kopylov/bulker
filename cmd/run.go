package cmd

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateRunCommand() *cobra.Command {
	var filter = runner.Filter{}

	var result = &cobra.Command{
		Use:   "run",
		Short: "Runs a custom command against each repository",
		Long: `Runs a custom command against each repository.
In order to pass the command to execute with its options, use "-- <command> [options]" notation, 
according to POSIX standard
See https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html#tag_12_02 :: Guideline 10
Example: "bulker run -- mvn -B -q clean"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("command to run is not found")
			}

			newRunner, err := runner.NewRunner(utils.GetConfiguredFS(), config.ReadConfig(), &filter)
			if err != nil {
				return err
			}

			err = newRunner.Run(
				func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
					return shell.RunCommand(runContext.Repo.Path, args[0], args[1:]...)
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	filter.AddCommandFlags(result)

	return result
}
