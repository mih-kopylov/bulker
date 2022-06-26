package gitops

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/spf13/afero"
	"os"
)

type CloneResult int

const (
	ClonedSuccessfully CloneResult = iota
	AlreadyCloned
	CloneError
)

func CloneRepo(fs afero.Fs, repo *model.Repo) (CloneResult, error) {
	_, err := fs.Stat(repo.Path)
	if err == nil {
		_, err := shell.RunCommand(repo.Path, "git", "status")
		if err == nil {
			return AlreadyCloned, nil
		}
		return CloneError, fmt.Errorf(
			"repository directory already exists, and it is not a repository: directory=%v",
			repo.Path,
		)
	}

	if !errors.Is(err, os.ErrNotExist) {
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
