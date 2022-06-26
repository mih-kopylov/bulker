package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/git"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Runs git operations on the repositories",
}

func init() {
	gitCmd.AddCommand(git.CloneCmd)
}
