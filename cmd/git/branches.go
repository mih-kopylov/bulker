package git

import (
	"github.com/mih-kopylov/bulker/cmd/git/branches"
	"github.com/spf13/cobra"
)

func CreateBranhesCommand() *cobra.Command {
	var result = &cobra.Command{
		Use:     "branches",
		Short:   "Manages git branches",
		Aliases: []string{"branch"},
	}

	result.AddCommand(branches.CreateListCommand())
	result.AddCommand(branches.CreateCheckoutCommand())

	return result
}
