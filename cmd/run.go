package cmd

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateRunCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}

	var result = &cobra.Command{
		Use:   "run",
		Short: "Runs a custom command against each repository",
		Long: `Runs a custom command against each repository.
In order to pass the command to execute with its options, use "-- <command> [options]" notation, 
according to POSIX standard

See https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html#tag_12_02 :: Guideline 10

Example: "bulker run -- mvn -B -q clean"`,
		Args: cobra.MinimumNArgs(1),
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				output, err := sh.RunCommand(runContext.Repo.Path, runContext.Args[0], runContext.Args[1:]...)
				if err != nil {
					return nil, fmt.Errorf("failed to run %v: %v %w", runContext.Args, output, err)
				}
				return output, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	return result
}
