package runner

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/spf13/cobra"
)

type Filter struct {
	Names []string
	Tags  []string
}

func (f Filter) Matches(_ settings.Repo) bool {
	return true
}

func (f Filter) AddCommandFlags(command *cobra.Command) {
	command.Flags().StringSliceVarP(&f.Names, "name", "n", []string{}, "Names of the repositories to process")
	command.Flags().StringSliceVarP(&f.Tags, "tag", "t", []string{}, "Tags of the repositories to process")
}
