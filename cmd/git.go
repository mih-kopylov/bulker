package cmd

import (
	"github.com/mih-kopylov/bulker/cmd/git"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/cobra"
)

func CreateGitCommand(sh shell.Shell) *cobra.Command {
	var result = &cobra.Command{
		Use:   "git",
		Short: "Runs git operations on the repositories",
	}

	result.AddCommand(git.CreateCloneCommand(sh))
	result.AddCommand(git.CreateFetchCommand(sh))
	result.AddCommand(git.CreatePullCommand(sh))
	result.AddCommand(git.CreatePushCommand(sh))
	result.AddCommand(git.CreateBranchesCommand(sh))
	result.AddCommand(git.CreateCommitCommand(sh))

	return result
}
