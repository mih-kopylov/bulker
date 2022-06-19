package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/repo"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Configures repositories that bulker works with",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	repoCmd.AddCommand(repo.ListCmd)
	repoCmd.AddCommand(repo.AddCmd)
	repoCmd.AddCommand(repo.RemoveCmd)
}
