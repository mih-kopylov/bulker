package gitops

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/afero"
	"os"
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

func Status(fs afero.Fs, repo *model.Repo) (StatusResult, error) {
	_, err := fs.Stat(repo.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return StatusMissing, nil
		} else {
			return StatusError, fmt.Errorf("failed to get stat of the directory %v: %w", repo.Path, err)
		}
	}

	statusResult, err := shell.RunCommand(repo.Path, "git", "status")
	if err != nil {
		return StatusError, err
	}

	if strings.Contains(statusResult, "working tree clean") {
		return StatusOk, nil
	}

	return StatusDirty, nil
}
