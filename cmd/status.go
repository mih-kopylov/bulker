package cmd

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/gitops"
	"github.com/mih-kopylov/bulker/internal/runner"
	"github.com/spf13/cobra"
	"strings"
)

func CreateStatusCommand() *cobra.Command {
	var filter = runner.Filter{}

	flags := struct {
		status StatusFilter
		ref    RefFilter
	}{}

	var result = &cobra.Command{
		Use:   "status",
		Short: "Prints status of all registered repositories",
		Long: `Prints status of all registered repositories. Status value can be one of:
* Clean - the repository successfully cloned, there are no uncommitted changes
* Dirty - the repository successfully cloned, but there are uncommitted changes
* Missing - the repository is not cloned yet`,
		RunE: runner.NewCommandRunner(
			&filter, func(ctx context.Context, runContext *runner.RunContext) (interface{}, error) {
				repoStatus, ref, err := gitops.Status(runContext.Repo)
				if err != nil {
					return nil, fmt.Errorf("failed to get status: %w", err)
				}

				type result struct {
					Status string
					Ref    string
				}

				if flags.status.Matches(repoStatus.String()) && flags.ref.Matches(ref) {
					return result{repoStatus.String(), ref}, nil
				}
				return nil, nil
			},
		),
	}

	filter.AddCommandFlags(result)

	result.Flags().Var(
		&flags.status, "status", fmt.Sprintf(
			`Keep repositories with specified status.
Examples: 
"bulker status --status clean" - will keep only repositories with "%v" status
"bulker status --status !clean" - will keep repositories with any status except "%v"`, gitops.StatusClean,
			gitops.StatusClean,
		),
	)
	result.Flags().Var(
		&flags.ref, "ref",
		`Keep repositories with specified ref.
Examples: 
"bulker status --ref master" - will keep only repositories with "master" ref
"bulker status --ref !master" - will keep repositories with any ref except "master"`,
	)

	return result
}

// StatusFilter implements Value in spf13/pflag for custom flag type
type StatusFilter string

func (f *StatusFilter) String() string {
	return string(*f)
}

func (f *StatusFilter) Set(v string) error {
	*f = StatusFilter(v)
	return nil
}

func (f *StatusFilter) Type() string {
	return "StatusFilter"
}

func (f *StatusFilter) Matches(status string) bool {
	if *f == "" {
		return true
	}

	negated, value := runner.ParseNegated(string(*f))

	if negated && !strings.EqualFold(value, status) {
		return true
	}

	if !negated && strings.EqualFold(value, status) {
		return true
	}

	return false
}

// RefFilter implements Value in spf13/pflag for custom flag type
type RefFilter string

func (f *RefFilter) String() string {
	return string(*f)
}

func (f *RefFilter) Set(v string) error {
	*f = RefFilter(v)
	return nil
}

func (f *RefFilter) Type() string {
	return "RefFilter"
}

func (f *RefFilter) Matches(ref string) bool {
	if *f == "" {
		return true
	}

	negated, value := runner.ParseNegated(string(*f))

	if negated && !strings.EqualFold(value, ref) {
		return true
	}

	if !negated && strings.EqualFold(value, ref) {
		return true
	}

	return false
}
