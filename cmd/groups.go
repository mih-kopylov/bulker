package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/groups"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateGroupsCommand(sh shell.Shell) *cobra.Command {
	var result = &cobra.Command{
		Use:     "groups",
		Short:   "Manages repositories groups",
		Aliases: []string{"group"},
	}

	result.AddCommand(groups.CreateListCommand(sh))
	result.AddCommand(groups.CreateGetCommand(sh))
	result.AddCommand(groups.CreateCreateCommand(sh))
	result.AddCommand(groups.CreateAppendCommand(sh))
	result.AddCommand(groups.CreateExcludeCommand(sh))
	result.AddCommand(groups.CreateRemoveCommand(sh))
	result.AddCommand(groups.CreateCleanCommand(sh))

	return result
}
