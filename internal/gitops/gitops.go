package gitops

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/afero"
	"os"
	"regexp"
	"strings"
)

type CloneResult int

const (
	ClonedSuccessfully CloneResult = iota
	ClonedAlready
	ClonedAgain
	CloneError
)

type StatusResult int

const (
	StatusOk StatusResult = iota
	StatusDirty
	StatusMissing
	StatusError
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

const (
	GitModeAll    GitMode = "all"
	GitModeLocal  GitMode = "local"
	GitModeRemote GitMode = "remote"
)

func CloneRepo(fs afero.Fs, repo *model.Repo, recreate bool) (CloneResult, error) {
	_, err := fs.Stat(repo.Path)
	if err == nil {
		_, err := shell.RunCommand(repo.Path, "git", "status")
		if err != nil {
			return CloneError, fmt.Errorf(
				"repository directory already exists, and it is not a repository: directory=%v",
				repo.Path,
			)
		}
		if !recreate {
			return ClonedAlready, nil
		}
		err = fs.RemoveAll(repo.Path)
		if err != nil {
			return CloneError, err
		}
	}

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return CloneError, fmt.Errorf("expected ErrNotExist but another found: %w", err)
	}

	err = fs.MkdirAll(repo.Path, 0700)
	if err != nil {
		return CloneError, fmt.Errorf("failed to create directory: directory=%v, error=%w", repo.Path, err)
	}

	_, err = shell.RunCommand(repo.Path, "git", "clone", repo.Url, ".")
	if err != nil {
		return CloneError, fmt.Errorf("failed to clone repository: repo=%v, error=%w", repo.Name, err)
	}

	return ClonedSuccessfully, nil
}

func Fetch(fs afero.Fs, repo *model.Repo) error {
	_, err := fs.Stat(repo.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("repository not cloned")
		}
		return err
	}

	_, err = shell.RunCommand(repo.Path, "git", "fetch")
	if err != nil {
		return fmt.Errorf("failed to fetch remote: %w", err)
	}

	return nil
}

func Pull(fs afero.Fs, repo *model.Repo) error {
	err := checkRepoExists(fs, repo)
	if err != nil {
		return err
	}

	_, err = shell.RunCommand(repo.Path, "git", "pull")
	if err != nil {
		return fmt.Errorf("failed to pull remote: %w", err)
	}

	return nil
}

func Status(fs afero.Fs, repo *model.Repo) (StatusResult, string, error) {
	_, err := fs.Stat(repo.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return StatusMissing, "", nil
		} else {
			return StatusError, "", fmt.Errorf("failed to get stat of the directory %v: %w", repo.Path, err)
		}
	}

	statusResult, err := shell.RunCommand(repo.Path, "git", "status")
	if err != nil {
		return StatusError, "", err
	}

	ref, err := parseHeadRef(statusResult)
	if err != nil {
		return StatusError, "", err
	}

	if strings.Contains(statusResult, "working tree clean") {
		return StatusOk, ref, nil
	}

	return StatusDirty, ref, nil
}

func GetBranches(fs afero.Fs, repo *model.Repo, mode GitMode, pattern string) ([]Branch, error) {
	err := checkRepoExists(fs, repo)
	if err != nil {
		return nil, err
	}

	reg, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return nil, err
	}

	outputString, err := shell.RunCommand(repo.Path, "git", "branch", "-a", "--format=%(refname)")
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	branches, err := parseBranches(outputString)
	if err != nil {
		return nil, err
	}

	var result []Branch
	for _, branch := range branches {
		if mode == GitModeLocal && !branch.IsLocal() {
			continue
		}

		if mode == GitModeRemote && branch.IsLocal() {
			continue
		}

		matched := reg.MatchString(branch.Name)
		if !matched {
			continue
		}

		result = append(result, branch)
	}

	return result, nil

}

func parseBranches(consoleOutputString string) ([]Branch, error) {
	var result []Branch
	for _, outputBranchName := range strings.Fields(consoleOutputString) {
		branch, err := parseBranch(outputBranchName)
		if err != nil {
			return nil, err
		}

		if branch.Name == Head {
			continue
		}
		result = append(result, *branch)
	}

	return result, nil
}

func checkRepoExists(fs afero.Fs, repo *model.Repo) error {
	_, err := fs.Stat(repo.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("repository not cloned")
		}
		return err
	}

	return nil
}

func parseHeadRef(statusResult string) (string, error) {
	reg, err := regexp.Compile("On branch (.+)\n.*")
	if err != nil {
		return "", err
	}

	submatch := reg.FindStringSubmatch(statusResult)
	if submatch != nil {
		return submatch[1], nil
	}

	reg, err = regexp.Compile("HEAD detached at (.+)\n.*")
	if err != nil {
		return "", err
	}

	submatch = reg.FindStringSubmatch(statusResult)
	if submatch != nil {
		return submatch[1], nil
	}

	return "", fmt.Errorf("can't parse status result for head reference: %v", statusResult)
}
