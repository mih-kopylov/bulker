package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/groups"
	"github.com/spf13/cobra"
)

func CreateGroupsCommand() *cobra.Command {
	var result = &cobra.Command{
		Use:   "groups",
		Short: "Manages repositories groups",
	}

	result.AddCommand(groups.CreateListCommand())
	result.AddCommand(groups.CreateGetCommand())
	result.AddCommand(groups.CreateCreateCommand())
	result.AddCommand(groups.CreateUpdateCommand())
	result.AddCommand(groups.CreateRemoveCommand())
	result.AddCommand(groups.CreateCleanCommand())

	return result
}
