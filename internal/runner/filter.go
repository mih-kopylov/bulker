package runner

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"strings"
)

type Filter struct {
	Names []string
	Tags  []string
}

var runMode = config.Parallel

func (f *Filter) Matches(repo settings.Repo) bool {
	return matchesName(repo.Name, f.Names) && matchesTags(repo.Tags, f.Tags)
}

func (f *Filter) AddCommandFlags(command *cobra.Command) {
	command.Flags().StringSliceVarP(&f.Names, "name", "n", []string{}, "Names of the repositories to process")
	command.Flags().StringSliceVarP(&f.Tags, "tag", "t", []string{}, "Tags of the repositories to process")

	// in order to viper read configuration from flag that is added multiple times (in different commands),
	// all the flags with the same name should have the same storage, which is a package variable
	command.PersistentFlags().Var(&runMode,
		"run-mode",
		"Parallel (par) or sequential (seq) run mode for repositories processing",
	)
	utils.BindFlag(command.PersistentFlags().Lookup("run-mode"), "runMode")
}

const negatePrefix = "-"

// matchesName repoName should match any of filterNames
func matchesName(repoName string, filterNames []string) bool {
	if len(filterNames) == 0 {
		return true
	}

	for _, filterName := range filterNames {
		if strings.HasPrefix(filterName, negatePrefix) && repoName != filterName[1:] {
			return true
		}
		if !strings.HasPrefix(filterName, negatePrefix) && repoName == filterName {
			return true
		}
	}

	return false
}

// matchesTags repoTags should match all filterTags
func matchesTags(repoTags []string, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}

	for _, filterTag := range filterTags {
		if strings.HasPrefix(filterTag, negatePrefix) && slices.Contains(repoTags, filterTag[1:]) {
			return false
		}

		if !strings.HasPrefix(filterTag, negatePrefix) && !slices.Contains(repoTags, filterTag) {
			return false
		}
	}

	return true
}
