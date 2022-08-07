package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/files"
	"github.com/spf13/cobra"
)

func CreateFilesCommand() *cobra.Command {
	var result = &cobra.Command{
		Use:     "files",
		Short:   "Manages files in the repositories",
		Aliases: []string{"file"},
	}

	result.AddCommand(files.CreateCopyCommand())

	return result
}
