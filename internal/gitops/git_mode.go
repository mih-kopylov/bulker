package gitops

import "fmt"

type GitMode string

func (g *GitMode) String() string {
	return string(*g)
}

func (g *GitMode) Set(v string) error {
	switch v {
	case string(GitModeAll), string(GitModeLocal), string(GitModeRemote):
		*g = GitMode(v)
		return nil
	default:
		return fmt.Errorf("must be one of '%s' '%s' '%s'", GitModeAll, GitModeLocal, GitModeRemote)
	}
}

func (g *GitMode) Type() string {
	return "GitMode"
}

func (g *GitMode) Includes(mode GitMode) bool {
	if *g == GitModeAll {
		return true
	}
	return *g == mode
}

const (
	GitModeAll    GitMode = "all"
	GitModeLocal  GitMode = "local"
	GitModeRemote GitMode = "remote"
)
