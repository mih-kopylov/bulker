package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/repos"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateReposCommand(sh shell.Shell) *cobra.Command {
	var result = &cobra.Command{
		Use:     "repos",
		Short:   "Configures repositories that bulker works with",
		Aliases: []string{"repo"},
	}

	result.AddCommand(repos.CreateListCommand(sh))
	result.AddCommand(repos.CreateAddCommand(sh))
	result.AddCommand(repos.CreateRemoveCommand(sh))
	result.AddCommand(repos.CreateExportCommand(sh))
	result.AddCommand(repos.CreateImportCommand(sh))

	return result
}
