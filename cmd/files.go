package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/files"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateFilesCommand(sh shell.Shell) *cobra.Command {
	var result = &cobra.Command{
		Use:     "files",
		Short:   "Manages files in the repositories",
		Aliases: []string{"file"},
	}

	result.AddCommand(files.CreateCopyCommand(sh))
	result.AddCommand(files.CreateRenameCommand(sh))
	result.AddCommand(files.CreateRemoveCommand(sh))
	result.AddCommand(files.CreateSearchCommand(sh))
	result.AddCommand(files.CreateReplaceCommand(sh))

	return result
}
