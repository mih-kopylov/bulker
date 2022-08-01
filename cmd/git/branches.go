package git

import (
	"github.com/mih-kopylov/bulker/cmd/git/branches"
	"github.com/spf13/cobra"
)

func CreateBranchesCommand() *cobra.Command {
	var result = &cobra.Command{
		Use:     "branches",
		Short:   "Manages git branches",
		Aliases: []string{"branch"},
	}

	result.AddCommand(branches.CreateListCommand())
	result.AddCommand(branches.CreateCheckoutCommand())
	result.AddCommand(branches.CreateCreateCommand())
	result.AddCommand(branches.CreateRemoveCommand())

	return result
}
