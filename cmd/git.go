package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/git"
	"github.com/spf13/cobra"
)

func CreateGitCommand() *cobra.Command {
	var result = &cobra.Command{
		Use:   "git",
		Short: "Runs git operations on the repositories",
	}

	result.AddCommand(git.CreateCloneCommand())
	result.AddCommand(git.CreateFetchCommand())
	result.AddCommand(git.CreatePullCommand())
	result.AddCommand(git.CreateBranhesCommand())

	return result
}
