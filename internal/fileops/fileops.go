package fileops

import (
	"errors"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/spf13/afero"
	"path/filepath"
)

var (
	ErrSourceNotFound      = errors.New("source not found")
	ErrTargetAlreadyExists = errors.New("target already exists")
)

func Copy(fs afero.Fs, repo *model.Repo, source string, target string, force bool) (string, string, error) {
	sourceAbs := filepath.Join(repo.Path, source)
	targetAbs := filepath.Join(repo.Path, target)

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
