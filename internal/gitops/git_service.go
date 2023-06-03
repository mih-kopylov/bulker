package gitops

import (
	"bytes"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
	"time"
)

type GitService struct {
	sh shell.Shell
}

func NewGitService(sh shell.Shell) GitService {
	return GitService{sh: sh}
}

func (g *GitService) CloneRepo(repo *model.Repo, recreate bool) (CloneResult, error) {
	exists, err := utils.Exists(repo.Path)
	if err != nil {
		return CloneError, err
	}

	emptyDir, _ := utils.EmptyDir(repo.Path)

	originalDirectoryDeleted := false

	if exists && !emptyDir {
		if recreate {
			err := os.RemoveAll(repo.Path)
			if err != nil {
				return CloneError, fmt.Errorf("failed to delete directory for recreation: %w", err)
			}

			exists = false
			emptyDir = false
			originalDirectoryDeleted = true
		} else {
			_, _, err := g.Status(repo)
			if err != nil {
				return CloneError, err
			}
			return ClonedAlready, nil
		}
	}

	if !exists {
		err = os.MkdirAll(repo.Path, 0700)
		if err != nil {
			return CloneError, fmt.Errorf("failed to create directory: directory=%v, error=%w", repo.Path, err)
		}
	}

	output, err := g.sh.RunCommand(repo.Path, "git", "clone", repo.Url, ".")
	if err != nil {
		return CloneError, fmt.Errorf("failed to clone repository: %v, %w", output, err)
	}

	if originalDirectoryDeleted {
		return ClonedAgain, nil
	}

	return ClonedSuccessfully, nil
}

func (g *GitService) Fetch(repo *model.Repo) error {
	output, err := g.sh.RunCommand(repo.Path, "git", "fetch", "--prune")
	if err != nil {
		return fmt.Errorf("failed to fetch remote: %v, %w", output, err)
	}

	return nil
}

func (g *GitService) Pull(repo *model.Repo) error {
	output, err := g.sh.RunCommand(repo.Path, "git", "pull", "--prune")
	if err != nil {
		if strings.Contains(output, "There is no tracking information for the current branch") {
			return fmt.Errorf("no remote upstream configured")
		}
		return fmt.Errorf("failed to pull remote: %v, %w", output, err)
	}

	return nil
}

func (g *GitService) Push(repo *model.Repo, branch string, allBranches bool, force bool) error {
	remote, err := g.getTheOnlyRemote(repo)
	if err != nil {
		return err
	}

	arguments := []string{"push", "--set-upstream", remote}
	if allBranches {
		if force {
			return errors.New("incompatible 'all' and 'force' modes")
		}
		arguments = append(arguments, "--all")
	} else {
		if branch == "" {
			return errors.New("incompatible 'branch' and 'allBranches' parameters values")
		}
		arguments = append(arguments, branch)
	}

	if force {
		arguments = append(arguments, "--force")
	}

	output, err := g.sh.RunCommand(repo.Path, "git", arguments...)
	if err != nil {
		return fmt.Errorf("failed to push to remote: %v, %w", output, err)
	}

	return nil
}

func (g *GitService) Status(repo *model.Repo) (StatusResult, string, error) {
	err := fileops.CheckRepoExists(repo)
	if err != nil {
		if errors.Is(err, fileops.ErrRepositoryNotCloned) {
			return StatusMissing, "", nil
		} else {
			return StatusError, "", fmt.Errorf("failed to get stat of the directory %v: %w", repo.Path, err)
		}
	}

	statusOutput, err := g.sh.RunCommand(repo.Path, "git", "status")
	if err != nil {
		return StatusError, "", fmt.Errorf("failed to get git status: %v, %w", statusOutput, err)
	}

	ref, err := g.parseHeadRef(statusOutput)
	if err != nil {
		return StatusError, "", err
	}

	if strings.Contains(statusOutput, "working tree clean") {
		return StatusClean, ref, nil
	}

	return StatusDirty, ref, nil
}

func (g *GitService) CreateBranch(repo *model.Repo, name string) (CreateResult, error) {
	branches, err := g.GetBranches(repo, GitModeAll, name)
	if err != nil {
		return CreateError, err
	}

	if len(branches) > 0 {
		return CreateError, fmt.Errorf("branch already exists")
	}

	output, err := g.sh.RunCommand(repo.Path, "git", "branch", name)
	if err != nil {
		return CreateError, fmt.Errorf("failed to create branch: %v, %w", output, err)
	}

	return CreateOk, nil
}

func (g *GitService) RemoveBranch(repo *model.Repo, name string, mode GitMode) (string, error) {
	branches, err := g.GetBranches(repo, mode, name)
	if err != nil {
		return "", err
	}

	if len(branches) == 0 {
		return "", fmt.Errorf("branch not found")
	}

	buffer := bytes.Buffer{}
	for _, branch := range branches {
		if branch.IsLocal() {
			output, err := g.sh.RunCommand(repo.Path, "git", "branch", "-D", name)
			if err != nil {
				if strings.Contains(output, "checked out at") {
					return "", fmt.Errorf("the branch is checked out")
				}
				return "", fmt.Errorf("failed to remove local branch: %v %w", output, err)
			}
			buffer.WriteString(fmt.Sprintf("%v: removed\n", branch.Short()))
		} else {
			output, err := g.sh.RunCommand(repo.Path, "git", "push", branch.Remote, "--delete", branch.Name)
			if err != nil {
				return "", fmt.Errorf("failed to remove remove branch: %v %w", output, err)
			}
			buffer.WriteString(fmt.Sprintf("%v: removed\n", branch.Short()))
		}
	}

	return strings.TrimSpace(buffer.String()), nil
}

func (g *GitService) CleanBranches(repo *model.Repo, mode GitMode) (string, error) {
	result := bytes.Buffer{}

	remote, err := g.getTheOnlyRemote(repo)
	if err != nil {
		return "", err
	}

	defaultRemoteBranch, err := g.getDefaultRemoteBranch(repo, remote)
	if err != nil {
		return "", err
	}

	if mode.Includes(GitModeLocal) {
		err := g.cleanLocalBranches(repo, defaultRemoteBranch, &result)
		if err != nil {
			return "", err
		}
	}
	if mode.Includes(GitModeRemote) {
		err := g.cleanRemoteBranches(repo, remote, defaultRemoteBranch, &result)
		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(result.String()), nil
}

func (g *GitService) Commit(repo *model.Repo, pattern string, message string) error {
	if pattern == "" {
		pattern = "**"
	}

	output, err := g.sh.RunCommand(repo.Path, "git", "add", pattern)
	if err != nil {
		return fmt.Errorf("failed to add changes to stage: %v %w", output, err)
	}

	output, err = g.sh.RunCommand(repo.Path, "git", "commit", "-m", message)
	if err != nil {
		return fmt.Errorf("failed to commit: %v %w", output, err)
	}

	return nil
}

func (g *GitService) Checkout(repo *model.Repo, ref string) (CheckoutResult, error) {
	branches, err := g.GetBranches(repo, GitModeAll, ref)
	if err != nil {
		return CheckoutError, err
	}

	if len(branches) == 0 {
		return CheckoutNotFound, nil
	}

	output, err := g.sh.RunCommand(repo.Path, "git", "checkout", ref)
	if err != nil {
		return CheckoutError, fmt.Errorf("failed to checkout: %v, %w", output, err)
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

func (g *GitService) Discard(repo *model.Repo) error {
	output, err := g.sh.RunCommand(repo.Path, "git", "reset", "--hard", "HEAD")
	if err != nil {
		return fmt.Errorf("failed to reset: %v, %w", output, err)
	}

	return nil
}

func (g *GitService) GetBranches(repo *model.Repo, mode GitMode, pattern string) ([]Branch, error) {
	reg, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return nil, err
	}

	output, err := g.sh.RunCommand(repo.Path, "git", "branch", "-a", "--format=%(refname)")
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %v, %w", output, err)
	}

	branches, err := g.parseBranches(output)
	if err != nil {
		return nil, err
	}

	var result []Branch
	for _, branch := range branches {
		if !mode.Includes(branch.GetGitMode()) {
			continue
		}

		matchedRegexp := reg.MatchString(branch.Name)
		// when searched by 'origin/my-branch' with intention to find the remote branch.
		// effectively the same as search by 'my-branch' with mode 'remote'
		matchedEquality := pattern == branch.Short()
		if !matchedRegexp && !matchedEquality {
			continue
		}

		result = append(result, branch)
	}

	return result, nil

}

func (g *GitService) GetDefaultBranch(repo *model.Repo) (*Branch, error) {
	remote, err := g.getTheOnlyRemote(repo)
	if err != nil {
		return nil, err
	}

	return g.getDefaultRemoteBranch(repo, remote)
}

// GetUnmergedBranches returns branches that are not merged to `ref`
func (g *GitService) GetUnmergedBranches(repo *model.Repo, mode GitMode, ref string) ([]Branch, error) {
	branches, err := g.getUnmergedBranches(repo, ref)
	if err != nil {
		return nil, err
	}

	if !mode.Includes(GitModeLocal) {
		branches = lo.Filter(
			branches, func(item Branch, _ int) bool {
				return !item.IsLocal()
			},
		)
	}
	if !mode.Includes(GitModeRemote) {
		branches = lo.Filter(
			branches, func(item Branch, _ int) bool {
				return item.IsLocal()
			},
		)
	}

	return branches, nil
}

// GetUnmergedCommits returns commits from branch `branch` that are not merged to ref `ref`.
// The commits are ordered by `committedAt` attribute
// Returns an error when `branch` and `ref` don't have a common parent
func (g *GitService) GetUnmergedCommits(repo *model.Repo, branch Branch, ref string) ([]Commit, error) {
	output, err := g.sh.RunCommand(
		repo.Path, "git", "--no-pager", "log", branch.Short(), "--not", ref,
		"--pretty=fuller",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get unmerged commits: %v, %w", output, err)
	}

	return g.parseCommits(output)
}

func (g *GitService) parseBranches(consoleOutputString string) ([]Branch, error) {
	var result []Branch
	outputBranchNames := strings.FieldsFunc(
		consoleOutputString, func(r rune) bool {
			return r == '\n'
		},
	)
	for _, outputBranchName := range outputBranchNames {
		branch, err := parseBranch(outputBranchName)
		if err != nil {
			if errors.Is(err, ErrDetachedHead) {
				continue
			}
			return nil, err
		}

		if branch.Name == Head {
			continue
		}
		result = append(result, *branch)
	}

	return result, nil
}

func (g *GitService) parseCommits(consoleOutputString string) ([]Commit, error) {
	commitRegexp := regexp.MustCompile(
		`commit\s+(\w+?)\n?.*
Author:\s+(.+?) <(.+?)>
AuthorDate:\s+(.+?)
Commit:\s+(.+?) <(.+?)>
CommitDate:\s+(.+?)
`,
	)

	var commits []Commit
	submatches := commitRegexp.FindAllStringSubmatch(consoleOutputString, -1)
	if submatches == nil {
		return nil, errors.New("failed to parse commits from console output")
	}
	for _, submatch := range submatches {
		commitId := submatch[1]
		authorName := submatch[2]
		authorEmail := submatch[3]
		authorDate := submatch[4]
		commitName := submatch[5]
		commitEmail := submatch[6]
		commitDate := submatch[7]

		authorDateTime, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", authorDate)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse author date from commit")
		}

		commitDateTime, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", commitDate)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse commit date from commit")
		}

		commit := Commit{
			Id:         commitId,
			AuthorDate: authorDateTime,
			Author: CommitUser{
				Name:  authorName,
				Email: authorEmail,
			},
			CommitDate: commitDateTime,
			Committer: CommitUser{
				Name:  commitName,
				Email: commitEmail,
			},
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

func (g *GitService) getDefaultRemoteBranch(repo *model.Repo, remote string) (*Branch, error) {
	output, err := g.sh.RunCommand(
		repo.Path, "git", "remote", "show", remote,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to show remote: %v, %w", output, err)
	}

	reg, err := regexp.Compile("HEAD branch: (.+)")
	if err != nil {
		return nil, err
	}

	var branchName string
	err = scanRegexp(output, reg, &branchName)
	if err != nil {
		return nil, err
	}

	return &Branch{Name: branchName, Remote: remote}, nil
}

func (g *GitService) parseHeadRef(statusResult string) (string, error) {
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

func (g *GitService) getUnmergedBranches(repo *model.Repo, ref string) ([]Branch, error) {
	output, err := g.sh.RunCommand(repo.Path, "git", "branch", "-a", "--format=%(refname)", "--no-merged", ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %v, %w", output, err)
	}

	return g.parseBranches(output)
}

func (g *GitService) cleanLocalBranches(repo *model.Repo, defaultRemoteBranch *Branch, result *bytes.Buffer) error {
	output, err := g.sh.RunCommand(
		repo.Path, "git", "branch", "-a", "--format=%(refname)", "--merged",
		defaultRemoteBranch.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to get branches: %v, %w", output, err)
	}

	branches, err := g.parseBranches(output)
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if !branch.IsLocal() {
			continue
		}
		if branch.Name == defaultRemoteBranch.Name {
			continue
		}
		output, err := g.sh.RunCommand(repo.Path, "git", "branch", "-d", branch.Name)
		if err != nil {
			if strings.Contains(output, "checked out at") {
				result.WriteString(fmt.Sprintf("%v: failed: %v\n", branch.Name, output))
				return fmt.Errorf("the branch is checked out")
			}
			result.WriteString(fmt.Sprintf("%v: failed: %v\n", branch.Name, output))
		} else {
			result.WriteString(fmt.Sprintf("%v: removed\n", branch.Name))
		}
	}

	return nil
}

func (g *GitService) cleanRemoteBranches(
	repo *model.Repo, remote string, defaultRemoteBranch *Branch, result *bytes.Buffer,
) error {
	output, err := g.sh.RunCommand(
		repo.Path, "git", "branch", "-a", "--format=%(refname)", "--merged",
		defaultRemoteBranch.Short(),
	)
	if err != nil {
		return fmt.Errorf("failed to get branches: %v, %w", output, err)
	}

	branches, err := g.parseBranches(output)
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if branch.IsLocal() {
			continue
		}
		if branch.Name == defaultRemoteBranch.Name {
			continue
		}
		output, err := g.sh.RunCommand(repo.Path, "git", "push", remote, "-d", branch.Name)
		if err != nil {
			result.WriteString(fmt.Sprintf("%v: failed: %v\n", branch.Short(), output))
		} else {
			result.WriteString(fmt.Sprintf("%v: removed\n", branch.Short()))
		}
	}

	return nil
}

func (g *GitService) getTheOnlyRemote(repo *model.Repo) (string, error) {
	output, err := g.sh.RunCommand(repo.Path, "git", "remote")
	if err != nil {
		return "", err
	}

	remotes := strings.Fields(output)
	if len(remotes) == 0 {
		return "", errors.New("no remotes found")
	}
	if len(remotes) > 1 {
		logrus.WithField("remotes", remotes).Warnf("expected to have only 1 remote, but %v found", len(remotes))
	}

	return remotes[0], nil
}
