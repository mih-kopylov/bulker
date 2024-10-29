package branches

import (
	"context"
	"maps"
	"slices"
	"sort"
	"time"

	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func CreateStaleCommand(sh shell.Shell) *cobra.Command {
	var filter = runner.Filter{}
	var flags struct {
		mode config.GitMode
		ref  string
		age  string
	}

	var result = &cobra.Command{
		Use:   "stale",
		Short: "Prints a list of stale repository branches",
		Long: `Prints a list of repository branches that are not merged to reference branch and have commits older
than specified age.`,
		RunE: runner.NewCommandRunnerForExistingRepos(
			&filter, sh, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {

				type staleBranch struct {
					Authors    []string
					LastCommit string
				}

				gitService := gitops.NewGitService(sh)

				latestCommitTime, err := utils.AgeToTime(&utils.RealClock{}, flags.age)
				if err != nil {
					return nil, err
				}

				if flags.ref == "" {
					defaultBranch, err := gitService.GetDefaultBranch(runContext.Repo)
					if err != nil {
						return nil, err
					}

					flags.ref = defaultBranch.Short()
				}

				unmergedBranches, err := gitService.GetUnmergedBranches(runContext.Repo, flags.mode, flags.ref)
				if err != nil {
					return nil, err
				}

				staleBranches := map[string]staleBranch{}

				for _, branch := range unmergedBranches {
					unmergedCommits, err := gitService.GetUnmergedCommits(runContext.Repo, branch, flags.ref)
					if err != nil {
						return nil, err
					}

					if len(unmergedCommits) == 0 {
						continue
					}

					lastCommit := lo.MaxBy(
						unmergedCommits, func(a gitops.Commit, b gitops.Commit) bool {
							return a.CommitDate.After(b.CommitDate)
						},
					)
					if !latestCommitTime.After(lastCommit.CommitDate) {
						continue
					}
					// was committed earlier than latestCommitTime, so too old, we consider the branch stale

					authors := getStaleBranchAuthors(unmergedCommits, 2)
					staleBranches[branch.Short()] = staleBranch{
						Authors:    authors,
						LastCommit: lastCommit.CommitDate.In(time.Local).Format(time.RFC3339),
					}
				}

				if len(staleBranches) == 0 {
					return nil, nil
				}

				bytes, err := yaml.Marshal(staleBranches)
				if err != nil {
					return nil, err
				}

				return string(bytes), nil
			},
		),
	}

	filter.AddCommandFlags(result)

	config.AddGitModeFlag(&flags.mode, result.Flags())
	result.Flags().StringVarP(
		&flags.ref, "ref", "r", "",
		`Git reference to compare branches with. If not set, the default repository branch is set`,
	)
	result.Flags().StringVarP(
		&flags.age, "age", "a", "",
		`Minimal age of the last commit in branches to consider them stale`,
	)
	utils.MarkFlagRequiredOrFail(result.Flags(), "age")

	return result
}

func getStaleBranchAuthors(commits []gitops.Commit, maxCount int) []string {
	authorCount := map[string]int{}
	for _, commit := range commits {
		authorCount[commit.Author.String()]++
	}

	authors := slices.Collect(maps.Keys(authorCount))
	sort.SliceStable(
		authors, func(i, j int) bool {
			return authorCount[authors[i]] > authorCount[authors[j]]
		},
	)
	if len(authors) <= maxCount {
		return authors
	}
	return authors[:maxCount]
}
