package repos

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateExportCommand() *cobra.Command {
	var flags struct {
		remote string
	}

	var result = &cobra.Command{
		Use:   "export",
		Short: "Exports the repositories configuration into an external git repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := settings.NewManager(config.ReadConfig())
			exportResult, err := manager.Export(flags.remote)
			if err != nil {
				return err
			}

			entityInfoMap := map[string]output.EntityInfo{}
			for repo, status := range exportResult {
				var statusString string
				switch status {
				case settings.ExportImportStatusCompleted:
					statusString = "exported"
				case settings.ExportImportStatusUpToDate:
					statusString = "up to date"
				default:
					statusString = fmt.Sprintf("status %v is not supported", status)
				}
				entityInfoMap[repo] = output.EntityInfo{Result: statusString}
			}

			err = output.Write("repo", entityInfoMap)
			if err != nil {
				return err
			}

			return nil
		},
	}

	result.Flags().StringVarP(&flags.remote, "remote", "r", "", "URL of the remote repository")
	utils.MarkFlagRequiredOrFail(result.Flags(), "remote")

	return result
}
