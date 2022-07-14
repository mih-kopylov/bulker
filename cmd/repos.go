package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/repos"
	"github.com/spf13/cobra"
)

func CreateReposCommand() *cobra.Command {
	var result = &cobra.Command{
		Use:   "repos",
		Short: "Configures repositories that bulker works with",
	}

	result.AddCommand(repos.CreateListCommand())
	result.AddCommand(repos.CreateAddCommand())
	result.AddCommand(repos.CreateRemoveCommand())
	result.AddCommand(repos.CreateExportCommand())

	return result
}
