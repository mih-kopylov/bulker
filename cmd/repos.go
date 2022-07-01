package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/repos"
	"github.com/spf13/cobra"
)

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Configures repositories that bulker works with",
}

func init() {
	reposCmd.AddCommand(repos.ListCmd)
	reposCmd.AddCommand(repos.AddCmd)
	reposCmd.AddCommand(repos.RemoveCmd)
}
