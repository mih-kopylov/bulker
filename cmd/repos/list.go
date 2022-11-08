package repos

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateListCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}

	var result = &cobra.Command{
		Use:   "list",
		Short: "Prints a list of supported repositories",
		RunE: runner.NewCommandRunner(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				return "", nil
			},
		),
	}

	filter.AddCommandFlags(result)

	return result
}
