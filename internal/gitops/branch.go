package gitops

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	RefPrefix       = "refs/"
	RefHeadPrefix   = RefPrefix + "heads/"
	RefRemotePrefix = RefPrefix + "remotes/"
	Head            = "HEAD"
)

var (
	ErrDetachedHead = errors.New("detached head")
)

type Branch struct {
	Name   string
	Remote string
}

func (b *Branch) String() string {
	if b.Remote == "" {
		return fmt.Sprintf("%s%s", RefHeadPrefix, b.Name)
	}
	return fmt.Sprintf("%s%s/%s", RefRemotePrefix, b.Remote, b.Name)
}

func (b *Branch) Short() string {
	if b.Remote == "" {
		return b.Name
	}
	return fmt.Sprintf("%s/%s", b.Remote, b.Name)
}

func (b *Branch) IsLocal() bool {
	return b.Remote == ""
}

func (b *Branch) GetGitMode() GitMode {
	if b.IsLocal() {
		return GitModeLocal
	}

	return GitModeRemote
}

func parseBranch(fullBranchName string) (*Branch, error) {
	branchName := ""
	branchRemote := ""

	if strings.HasPrefix(fullBranchName, "(HEAD detached at") {
		return nil, ErrDetachedHead
	} else if strings.HasPrefix(fullBranchName, RefHeadPrefix) {
		reg, err := regexp.Compile(RefHeadPrefix + "(.+)")
		if err != nil {
			return nil, err
		}

		err = scanRegexp(fullBranchName, reg, &branchName)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(fullBranchName, RefRemotePrefix) {
		reg, err := regexp.Compile(RefRemotePrefix + "(.+?)/(.+)")
		if err != nil {
			return nil, err
		}

		err = scanRegexp(fullBranchName, reg, &branchRemote, &branchName)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported branch name: %v", fullBranchName)
	}

	return &Branch{branchName, branchRemote}, nil
}

func scanRegexp(value string, reg *regexp.Regexp, groups ...*string) error {
	matchedGroups := reg.FindStringSubmatch(value)
	if len(matchedGroups) != len(groups)+1 {
		return fmt.Errorf("unexpected groups matched: groups=%v, regex=%v, value=%v", matchedGroups, reg, value)
	}

	for i, group := range groups {
		*group = matchedGroups[i+1]
	}

	return nil
}
