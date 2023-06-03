package git

import (
	"github.com/mih-kopylov/bulker/cmd/git/branches"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateBranchesCommand(sh shell.Shell) *cobra.Command {
	var result = &cobra.Command{
		Use:     "branches",
		Short:   "Manages git branches",
		Aliases: []string{"branch"},
	}

	result.AddCommand(branches.CreateListCommand(sh))
	result.AddCommand(branches.CreateCheckoutCommand(sh))
	result.AddCommand(branches.CreateCreateCommand(sh))
	result.AddCommand(branches.CreateRemoveCommand(sh))
	result.AddCommand(branches.CreateCleanCommand(sh))
	result.AddCommand(branches.CreateStaleCommand(sh))

	return result
}
