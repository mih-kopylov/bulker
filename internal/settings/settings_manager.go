package settings

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
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

	fileContent, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	err = afero.WriteFile(sm.fs, settingsFileName, fileContent, os.ModePerm)
	return err
}

func (sm *Manager) Export(remoteRepoUrl string) (map[string]ExportStatus, error) {
	logrus.WithField("remote", remoteRepoUrl).Debug("exporting repositories")
	settings, err := sm.Read()
	if err != nil {
		return nil, err
	}

	repoDir, err := os.MkdirTemp("", "bulker_remote_repo_*")
	if err != nil {
		return nil, err
	}
	defer func() {
		err := os.RemoveAll(repoDir)
		if err != nil {
			logrus.Warnf("can't remove temp directory %v: %v", repoDir, err)
		}
	}()

	_, err = shell.RunCommand(repoDir, "git", "clone", remoteRepoUrl, ".")
	if err != nil {
		return nil, err
	}

	modelToExport := fromSettings(settings)
	jsonBytes, err := yaml.Marshal(modelToExport)
	if err != nil {
		return nil, err
	}

	exportFileName := path.Join(repoDir, "repos.yaml")
	existingModel, err := readExistingModel(exportFileName)
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
		result := map[string]ExportStatus{}
		for _, repo := range settings.Repos {
			result[repo.Name] = ExportStatusUpToDate
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

	result := prepareResult(existingModel, modelToExport)

	return result, nil
}

func prepareResult(existingModel *exportModel, modelToExport *exportModel) map[string]ExportStatus {
	result := map[string]ExportStatus{}
	for repoName, repo := range modelToExport.Data.Repos {
		if existingModel.Version != modelToExport.Version {
			result[repoName] = ExportStatusExported
		} else {
			existingRepo := modelDataV1Repo{}
			for existingRepoIterName, existingRepoIter := range existingModel.Data.Repos {
				if existingRepoIterName == repoName {
					existingRepo = existingRepoIter
				}
			}
			if existingRepo.Equals(repo) {
				result[repoName] = ExportStatusUpToDate
			} else {
				result[repoName] = ExportStatusExported
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
