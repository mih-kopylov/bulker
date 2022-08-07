package fileops

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

var (
	ErrSourceNotFound      = errors.New("source not found")
	ErrTargetAlreadyExists = errors.New("target already exists")
)

func Copy(fs afero.Fs, repo *model.Repo, source string, target string, force bool) (string, string, error) {
	sourceAbs := filepath.Join(repo.Path, source)
	targetAbs := filepath.Join(repo.Path, target)

	err := CheckRepoExists(fs, repo)
	if err != nil {
		return sourceAbs, targetAbs, err
	}

	sourceExists, err := afero.Exists(fs, sourceAbs)
	if err != nil {
		return sourceAbs, targetAbs, err
	}
	if !sourceExists {
		return sourceAbs, targetAbs, ErrSourceNotFound
	}

	targetExists, err := afero.Exists(fs, targetAbs)
	if err != nil {
		return sourceAbs, targetAbs, err
	}
	if targetExists && !force {
		return sourceAbs, targetAbs, ErrTargetAlreadyExists
	}

	fileContent, err := afero.ReadFile(fs, sourceAbs)
	if err != nil {
		return sourceAbs, targetAbs, err
	}

	sourceFileInfo, err := fs.Stat(sourceAbs)
	if err != nil {
		return sourceAbs, targetAbs, err
	}

	err = afero.WriteFile(fs, targetAbs, fileContent, sourceFileInfo.Mode())
	if err != nil {
		return sourceAbs, targetAbs, err
	}

	return sourceAbs, targetAbs, nil
}

func Rename(fs afero.Fs, repo *model.Repo, source string, target string, force bool) (string, string, error) {
	sourceAbs := filepath.Join(repo.Path, source)
	targetAbs := filepath.Join(repo.Path, target)

	err := CheckRepoExists(fs, repo)
	if err != nil {
		return sourceAbs, targetAbs, err
	}

	sourceExists, err := afero.Exists(fs, sourceAbs)
	if err != nil {
		return sourceAbs, targetAbs, err
	}
	if !sourceExists {
		return sourceAbs, targetAbs, ErrSourceNotFound
	}

	targetExists, err := afero.Exists(fs, targetAbs)
	if err != nil {
		return sourceAbs, targetAbs, err
	}
	if targetExists && !force {
		return sourceAbs, targetAbs, ErrTargetAlreadyExists
	}

	err = fs.Rename(sourceAbs, targetAbs)
	if err != nil {
		return sourceAbs, targetAbs, err
	}

	return sourceAbs, targetAbs, nil
}

func Remove(fs afero.Fs, repo *model.Repo, pattern string) ([]string, error) {
	err := CheckRepoExists(fs, repo)
	if err != nil {
		return nil, err
	}

	matches, err := afero.Glob(fs, filepath.Join(repo.Path, pattern))
	if err != nil {
		return nil, err
	}

	var result []string
	for _, fileToRemove := range matches {
		err := fs.Remove(fileToRemove)
		if err != nil {
			result = append(result, fmt.Sprintf("%v: failed: %v", fileToRemove, err))
		} else {
			result = append(result, fmt.Sprintf("%v: removed", fileToRemove))
		}
	}

	return result, nil
}

var ErrRepositoryNotCloned = errors.New("repository not cloned")

func CheckRepoExists(fs afero.Fs, repo *model.Repo) error {
	_, err := fs.Stat(repo.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrRepositoryNotCloned
		}
		return err
	}

	return nil
}
