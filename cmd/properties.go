package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/props"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreatePropertiesCommand(sh shell.Shell) *cobra.Command {
	var result = &cobra.Command{
		Use:     "properties",
		Short:   "Manages configuration properties in different file formats",
		Aliases: []string{"property", "props", "prop"},
	}

	result.AddCommand(props.CreateGetCommand(sh))

	return result
}
