package groups

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateCreateCommand(sh shell.Shell) *cobra.Command {
	flags := struct {
		group string
		repos []string
		force bool
		from  string
	}{}

	var result = &cobra.Command{
		Use:   "create",
		Short: "Creates a new group with provided content",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(config.ReadConfig(), sh)

			sets, err := settingsManager.Read()
			if err != nil {
				return err
			}

			group, err := sets.AddGroup(flags.group)
			if err != nil {
				if !errors.Is(err, settings.ErrGroupAlreadyExists) || !flags.force {
					return fmt.Errorf("group already exists, use --force to recreate")
				}

				err = sets.RemoveGroup(flags.group)
				if err != nil {
					return err
				}

				group, err = sets.AddGroup(flags.group)
				if err != nil {
					return err
				}
			}

			entityInfoMap := map[string]output.EntityInfo{}

			repos, err := utils.GetReposFromStdInOrDefault(flags.repos)
			if err != nil {
				return err
			}

			for _, repoName := range repos {
				err := sets.AddRepoToGroup(group, repoName)
				if err != nil {
					entityInfoMap[repoName] = output.EntityInfo{Result: nil, Error: err}
				} else {
					entityInfoMap[repoName] = output.EntityInfo{Result: "created", Error: nil}
				}
			}

			err = settingsManager.Write(sets)
			if err != nil {
				return err
			}

			err = output.Write(cmd.OutOrStdout(), "repo", entityInfoMap)
			if err != nil {
				return err
			}

			return nil
		},
	}

	result.Flags().StringVarP(&flags.group, "group", "g", "", "Name of the group to remove")
	utils.MarkFlagRequiredOrFail(result.Flags(), "group")

	result.Flags().StringSliceVarP(
		&flags.repos, "name", "n", []string{}, "Names of the repositories to add to the group",
	)

	result.Flags().BoolVarP(
		&flags.force, "force", "f", false, "Recreate the group if a group with such a name already exists",
	)

	result.Flags().StringVar(
		&flags.from, "from", "", `Name of the group to use as a source for the new one. 
Use 'previous' to refer to the previous command result`,
	)
	result.MarkFlagsMutuallyExclusive("name", "from")

	utils.AddReadFromStdInFlag(result, "repo")

	return result
}
