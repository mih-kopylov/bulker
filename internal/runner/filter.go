package runner

import (
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
)

type Filter struct {
	Names  []string
	Tags   []string
	Groups []string
}

func (f *Filter) MatchesRepo(repo settings.Repo, groups []settings.Group) bool {
	return f.matchesName(repo.Name) &&
		f.matchesTags(repo.Tags) &&
		f.matchesGroups(repo.Name, groups)
}

func (f *Filter) FilterMatchingRepos(repos []settings.Repo, groups []settings.Group) []settings.Repo {
	var result []settings.Repo
	for _, repo := range repos {
		if !f.MatchesRepo(repo, groups) {
			continue
		}
		result = append(result, repo)
	}
	return result
}

func (f *Filter) AddCommandFlags(command *cobra.Command) {
	command.Flags().StringSliceVarP(
		&f.Names, "name", "n", []string{},
		"Names of the repositories to process. Can be regexp",
	)
	command.Flags().StringSliceVarP(&f.Tags, "tag", "t", []string{}, "Tags of the repositories to process")
	command.Flags().StringSliceVarP(&f.Groups, "group", "g", []string{}, "Groups of the repositories to process")
}

const negatePrefix = "!"

// ParseNegated gets a string value and analyzes if it is negated or not - effectively,
// if it starts with a prefix or not.
//
// It returns a bool flag and a clean value, without the prefix even for negated value
func ParseNegated(value string) (bool, string) {
	negated := false
	if strings.HasPrefix(value, negatePrefix) {
		negated = true
		value = value[len(negatePrefix):]
	}

	return negated, value
}

// matchesName repoName should match any of filterNames
func (f *Filter) matchesName(repoName string) bool {
	if len(f.Names) == 0 {
		return true
	}

	for _, filterName := range f.Names {
		negated, filterName := ParseNegated(filterName)

		matched, _ := regexp.MatchString("^"+filterName+"$", repoName)
		if negated {
			matched = !matched
		}
		if matched {
			return true
		}
	}

	return false
}

// matchesTags repoTags should match all filterTags
func (f *Filter) matchesTags(repoTags []string) bool {
	if len(f.Tags) == 0 {
		return true
	}

	for _, filterTag := range f.Tags {
		negated, filterTag := ParseNegated(filterTag)

		if negated && slices.Contains(repoTags, filterTag) {
			return false
		}

		if !negated && !slices.Contains(repoTags, filterTag) {
			return false
		}
	}

	return true
}

// matchesGroups repoName should match all group filters
func (f *Filter) matchesGroups(repoName string, allGroups []settings.Group) bool {
	if len(f.Groups) == 0 {
		return true
	}

	for _, filterGroupName := range f.Groups {
		negated, filterGroupName := ParseNegated(filterGroupName)

		filterGroupIndex := slices.IndexFunc(
			allGroups, func(group settings.Group) bool {
				return group.Name == filterGroupName
			},
		)
		if filterGroupIndex < 0 {
			// passed group name that is not found in settings
			if negated {
				continue
			}
			return false
		}
		filterGroup := allGroups[filterGroupIndex]
		contains := slices.Contains(filterGroup.Repos, repoName)
		matches := contains
		if negated {
			matches = !contains
		}
		if !matches {
			return false
		}
	}
	return true
}
