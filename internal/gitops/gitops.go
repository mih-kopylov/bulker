package gitops

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/afero"
	"golang.org/x/exp/slices"
	"os"
	"regexp"
	"strings"
)

var ErrRepositoryNotCloned = errors.New("repository not cloned")

type CloneResult int

const (
	ClonedSuccessfully CloneResult = iota
	ClonedAlready
	ClonedAgain
	CloneError
)

type StatusResult string

const (
	StatusClean   StatusResult = "Clean"
	StatusDirty   StatusResult = "Dirty"
	StatusMissing StatusResult = "Missing"
	StatusError   StatusResult = "Error"
)

func (r *StatusResult) String() string {
	return string(*r)
}

type CheckoutResult string

const (
	CheckoutOk       = "Success"
	CheckoutNotFound = "Not Found"
	CheckoutError    = "Error"
)

func (r *CheckoutResult) String() string {
	return string(*r)
}

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

	output, err := shell.RunCommand(repo.Path, "git", "pull")
	if err != nil {
		if strings.Contains(output, "There is no tracking information for the current branch") {
			return fmt.Errorf("no remote upstream configured")
		}
		return fmt.Errorf("failed to pull remote: %v %w", output, err)
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
		return StatusClean, ref, nil
	}

	return StatusDirty, ref, nil
}

func Checkout(fs afero.Fs, repo *model.Repo, ref string) (CheckoutResult, error) {
	err := checkRepoExists(fs, repo)
	if err != nil {
		return CheckoutError, err
	}

	branches, err := GetBranches(fs, repo, GitModeAll, ".*")
	if err != nil {
		return CheckoutError, err
	}

	branchIndex := slices.IndexFunc(
		branches, func(b Branch) bool {
			return b.Name == ref
		},
	)
	if branchIndex < 0 {
		return CheckoutNotFound, nil
	}

	output, err := shell.RunCommand(repo.Path, "git", "checkout", ref)
	if err != nil {
		return CheckoutError, err
	}

	if strings.Contains(output, "Already on") {
		return CheckoutOk, nil
	}

	if strings.Contains(output, "Switched to branch") {
		return CheckoutOk, nil
	}

	if strings.Contains(output, "Switched to a new branch") {
		return CheckoutOk, nil
	}

	return CheckoutError, fmt.Errorf("unknown checkout status: %v", output)
}

func Discard(fs afero.Fs, repo *model.Repo) error {
	err := checkRepoExists(fs, repo)
	if err != nil {
		return err
	}

	output, err := shell.RunCommand(repo.Path, "git", "reset", "--hard", "HEAD")
	if err != nil {
		return fmt.Errorf("failed to reset: %v %w", output, err)
	}

	return nil
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
			return ErrRepositoryNotCloned
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
