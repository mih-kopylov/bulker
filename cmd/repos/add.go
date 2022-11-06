package repos

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateAddCommand() *cobra.Command {
	var flags struct {
		name string
		url  string
		tags []string
	}

	var result = &cobra.Command{
		Use:   "add",
		Short: "Adds a new repository to the supported list",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsManager := settings.NewManager(config.ReadConfig())

			sets, err := settingsManager.Read()
			if err != nil {
				return err
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
				"repo",
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

	result.Flags().StringVar(&flags.name, "name", "", "Name of the repository")
	utils.MarkFlagRequiredOrFail(result.Flags(), "name")

	result.Flags().StringVar(&flags.url, "url", "", "URL of the repository")
	utils.MarkFlagRequiredOrFail(result.Flags(), "url")

	result.Flags().StringSliceVar(&flags.tags, "tags", []string{}, "Tags of the repository")

	return result
}
