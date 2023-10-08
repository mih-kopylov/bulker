package repos

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
	"path/filepath"
)

func CreateAddCommand(sh shell.Shell) *cobra.Command {
	var flags struct {
		name string
		url  string
		tags []string
	}

	var result = &cobra.Command{
		Use:   "add",
		Short: "Adds a new repository to the supported list",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(config.ReadConfig(), sh)

			sets, err := settingsManager.Read()
			if err != nil {
				return err
			}

			if flags.name == "" {
				flags.name = filepath.Base(flags.url)
				if filepath.Ext(flags.name) == ".git" {
					flags.name = flags.name[:len(flags.name)-len(".git")]
				}
			}

			err = sets.AddRepo(flags.name, flags.url, flags.tags)
			if err != nil {
				return err
			}

			err = settingsManager.Write(sets)
			if err != nil {
				return err
			}

			err = output.Write(
				cmd.OutOrStdout(), "repo",
				map[string]output.EntityInfo{
					flags.name: {Result: "added"},
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	result.Flags().StringVarP(
		&flags.name, "name", "n", "", `Name of the repository.
Use the name to override the directory name where the sources will be stored in the projects directory.
By default the name will be taken from the git repo name`,
	)

	result.Flags().StringVarP(
		&flags.url, "url", "u", "", `URL of the repository.
Effectively, the "origin" remote, the default one.
It will be used to clone the repository.`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "url")

	result.Flags().StringSliceVar(&flags.tags, "tags", []string{}, "Tags of the repository")

	return result
}
