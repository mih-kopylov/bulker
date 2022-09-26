package fileops

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrSourceNotFound      = errors.New("source not found")
	ErrTargetAlreadyExists = errors.New("target already exists")
	ErrRepositoryNotCloned = errors.New("repository not cloned")
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

type FileSearchResult struct {
	FileName string
	Matches  []string
}

func SearchFiles(
	fs afero.Fs, repo *model.Repo, pattern string, contains string, before int, after int,
) ([]FileSearchResult, error) {
	err := CheckRepoExists(fs, repo)
	if err != nil {
		return nil, err
	}

	matchedFiles, err := afero.Glob(fs, filepath.Join(repo.Path, pattern))
	if err != nil {
		return nil, err
	}

	var containsReg *regexp.Regexp
	if contains != "" {
		containsReg, err = regexp.Compile(contains)
		if err != nil {
			return nil, err
		}
	}

	var result []FileSearchResult
	for _, matchedFile := range matchedFiles {
		if containsReg == nil {
			result = append(
				result, FileSearchResult{
					FileName: matchedFile,
					Matches:  nil,
				},
			)
		} else {
			searchResult, err := SearchInFile(fs, matchedFile, containsReg, before, after)
			if err != nil {
				return nil, err
			}

			if searchResult != nil {
				result = append(result, *searchResult)
			}
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func SearchInFile(fs afero.Fs, fileName string, containsReg *regexp.Regexp, before int, after int) (
	*FileSearchResult, error,
) {
	if containsReg == nil {
		return nil, errors.New("containsReg is expected to be passed")
	}

	fileBytes, err := afero.ReadFile(fs, fileName)
	if err != nil {
		return nil, err
	}

	findings := containsReg.FindAllIndex(fileBytes, -1)
	if len(findings) == 0 {
		return nil, nil
	}
	result := &FileSearchResult{
		FileName: fileName,
		Matches:  []string{},
	}
	for _, finding := range findings {
		beforeFindingString := string(fileBytes[:finding[0]])
		foundLineNumber := strings.Count(beforeFindingString, "\n")

		fileContentLines := strings.Split(string(fileBytes), "\n")

		var foundResult []string

		for i := 0; i < before; i++ {
			lineIndex := foundLineNumber - before + i
			if lineIndex >= 0 {
				foundResult = append(foundResult, fileContentLines[lineIndex])
			}
		}
		foundResult = append(foundResult, fileContentLines[foundLineNumber])
		for i := 0; i < after; i++ {
			lineIndex := foundLineNumber + i + 1
			if lineIndex < len(fileContentLines) {
				foundResult = append(foundResult, fileContentLines[lineIndex])
			}
		}

		result.Matches = append(result.Matches, strings.Join(foundResult, "\n"))
	}
	return result, nil
}
