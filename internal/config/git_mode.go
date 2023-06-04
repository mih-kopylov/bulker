package config

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/spf13/pflag"
)

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

func AddGitModeFlag(storage *GitMode, flagSet *pflag.FlagSet) {
	defaultMode := ReadConfig().GitMode
	if defaultMode == "" {
		defaultMode = GitModeLocal
	}
	*storage = defaultMode

	flagSet.VarP(
		storage, "mode", "m", fmt.Sprintf(
			"Git work mode. Whether to process local repository or remote or both. Available types are: %s, %s, %s",
			GitModeAll, GitModeLocal, GitModeRemote,
		),
	)
	utils.BindFlag(flagSet.Lookup("mode"), "gitMode")
}
