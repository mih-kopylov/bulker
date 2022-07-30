package settings

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Manager struct {
	fs   afero.Fs
	conf *config.Config
}

func NewManager(fs afero.Fs, conf *config.Config) *Manager {
	return &Manager{
		fs:   fs,
		conf: conf,
	}
}

func (sm *Manager) Read() (*Settings, error) {
	settingsFileName := sm.conf.SettingsFileName

	settingsDirectory := filepath.Dir(settingsFileName)
	err := sm.fs.MkdirAll(settingsDirectory, 0777)
	if err != nil {
		return nil, err
	}

	fileContent, err := afero.ReadFile(sm.fs, settingsFileName)
	if err != nil {
		newInstance := &Settings{}
		err := sm.Write(newInstance)
		if err != nil {
			return nil, err
		}

		return newInstance, nil
	}

	result := &Settings{}
	err = yaml.Unmarshal(fileContent, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (sm *Manager) Write(settings *Settings) error {
	settingsFileName := sm.conf.SettingsFileName

	// make sure all data is sorted alphabetically
	slices.SortFunc(
		settings.Repos, func(a Repo, b Repo) bool {
			return strings.Compare(a.Name, b.Name) < 0
		},
	)
	slices.SortFunc(
		settings.Groups, func(a Group, b Group) bool {
			return strings.Compare(a.Name, b.Name) < 0
		},
	)
	for _, group := range settings.Groups {
		slices.Sort(group.Repos)
	}

	fileContent, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	err = afero.WriteFile(sm.fs, settingsFileName, fileContent, os.ModePerm)
	return err
}

func (sm *Manager) Export(remoteRepoUrl string) (map[string]ExportImportStatus, error) {
	logrus.WithField("remote", remoteRepoUrl).Debug("exporting repositories")
	settings, err := sm.Read()
	if err != nil {
		return nil, err
	}

	repoDir, cleanupFunc, err := sm.cloneRepo(remoteRepoUrl)
	if cleanupFunc != nil {
		defer cleanupFunc()
	}
	if err != nil {
		return nil, err
	}

	settingsModel := fromSettings(settings)
	jsonBytes, err := yaml.Marshal(settingsModel)
	if err != nil {
		return nil, err
	}

	exportFileName := filepath.Join(repoDir, exportImportFileName)
	fileModel, err := readExistingModel(exportFileName)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(exportFileName, jsonBytes, 0777)
	if err != nil {
		return nil, err
	}

	_, err = shell.RunCommand(repoDir, "git", "add", ".")
	if err != nil {
		return nil, err
	}

	statusOutput, err := shell.RunCommand(repoDir, "git", "status")
	if err != nil {
		return nil, err
	}

	if strings.Contains(statusOutput, "nothing to commit") {
		result := map[string]ExportImportStatus{}
		for _, repo := range settings.Repos {
			result[repo.Name] = ExportImportStatusUpToDate
		}
		return result, nil
	}

	_, err = shell.RunCommand(repoDir, "git", "commit", "-m", "Export bulker repositories")
	if err != nil {
		return nil, err
	}

	_, err = shell.RunCommand(repoDir, "git", "push")
	if err != nil {
		return nil, err
	}

	result := prepareResult(fileModel, settingsModel)

	return result, nil
}

func (sm *Manager) Import(remoteRepoUrl string) (map[string]ExportImportStatus, error) {
	logrus.WithField("remote", remoteRepoUrl).Debug("importing repositories")
	settings, err := sm.Read()
	if err != nil {
		return nil, err
	}

	repoDir, cleanupFunc, err := sm.cloneRepo(remoteRepoUrl)
	if cleanupFunc != nil {
		defer cleanupFunc()
	}
	if err != nil {
		return nil, err
	}

	settingsModel := fromSettings(settings)
	importFileName := filepath.Join(repoDir, exportImportFileName)
	fileModel, err := readExistingModel(importFileName)
	if err != nil {
		return nil, err
	}

	importedSettings, err := toSettings(fileModel)
	if err != nil {
		return nil, err
	}

	err = sm.Write(importedSettings)
	if err != nil {
		return nil, err
	}
	result := prepareResult(settingsModel, fileModel)

	return result, nil
}

const exportImportFileName = "repos.yaml"

func (sm *Manager) cloneRepo(remoteRepoUrl string) (repoDir string, cleanupFunc func(), err error) {
	repoDir, err = os.MkdirTemp("", "bulker_remote_repo_*")
	if err != nil {
		return "", nil, err
	}
	cleanupFunc = func() {
		err := os.RemoveAll(repoDir)
		if err != nil {
			logrus.Warnf("can't remove temp directory %v: %v", repoDir, err)
		}
	}

	_, err = shell.RunCommand(repoDir, "git", "clone", remoteRepoUrl, ".")
	if err != nil {
		return "", cleanupFunc, nil
	}
	return repoDir, cleanupFunc, err
}

func prepareResult(previousModel *exportModel, newModel *exportModel) map[string]ExportImportStatus {
	result := map[string]ExportImportStatus{}
	for repoName, repo := range newModel.Data.Repos {
		if previousModel.Version != newModel.Version {
			result[repoName] = ExportImportStatusCompleted
		} else {
			existingRepo := modelDataV1Repo{}
			for existingRepoIterName, existingRepoIter := range previousModel.Data.Repos {
				if existingRepoIterName == repoName {
					existingRepo = existingRepoIter
				}
			}
			if existingRepo.Equals(repo) {
				result[repoName] = ExportImportStatusUpToDate
			} else {
				result[repoName] = ExportImportStatusCompleted
			}
		}
	}
	return result
}

func readExistingModel(exportFileName string) (*exportModel, error) {
	fileBytes, err := ioutil.ReadFile(exportFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &exportModel{1, modelDataV1{Repos: map[string]modelDataV1Repo{}}}, nil
		}

		return nil, err
	}

	var readVersionModel map[string]any
	err = yaml.Unmarshal(fileBytes, &readVersionModel)
	if err != nil {
		return nil, err
	}

	version := readVersionModel["version"].(int)

	data := readVersionModel["data"]
	dataBytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	switch version {
	case 1:
		var resultData modelDataV1
		err := yaml.Unmarshal(dataBytes, &resultData)
		if err != nil {
			return nil, err
		}

		return &exportModel{Version: version, Data: resultData}, nil
	default:
		return nil, fmt.Errorf("version %v is not supported", version)
	}
}

func fromSettings(settings *Settings) *exportModel {
	data := modelDataV1{}
	data.Repos = map[string]modelDataV1Repo{}
	for _, repo := range settings.Repos {
		data.Repos[repo.Name] = modelDataV1Repo{
			Url:  repo.Url,
			Tags: repo.Tags,
		}
	}

	return &exportModel{1, data}
}

func toSettings(em *exportModel) (*Settings, error) {
	result := Settings{[]Repo{}, []Group{}}
	for repoName, repoData := range em.Data.Repos {
		err := result.AddRepo(repoName, repoData.Url, repoData.Tags)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}
