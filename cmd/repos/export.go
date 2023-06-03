package repos

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
)

func CreateExportCommand(sh shell.Shell) *cobra.Command {
	var flags struct {
		remote string
	}

	var result = &cobra.Command{
		Use:   "export",
		Short: "Exports the repositories configuration into an external git repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := settings.NewManager(config.ReadConfig(), sh)
			exportResult, err := manager.Export(flags.remote)
			if err != nil {
				return err
			}

			entityInfoMap := map[string]output.EntityInfo{}
			for repo, status := range exportResult {
				switch status {
				case settings.ExportImportStatusUpToDate:
					//omit in the result
				case settings.ExportImportStatusAdded:
					entityInfoMap[repo] = output.EntityInfo{Result: "exported"}
				case settings.ExportImportStatusRemoved:
					entityInfoMap[repo] = output.EntityInfo{Result: "removal exported"}
				default:
					entityInfoMap[repo] = output.EntityInfo{Error: fmt.Errorf("status %v is not supported", status)}
				}
			}

			err = output.Write(cmd.OutOrStdout(), "repo", entityInfoMap)
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
