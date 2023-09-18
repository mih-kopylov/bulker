package settings

import (
	"errors"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

const currentVersion = 1

type Manager struct {
	conf *config.Config
	sh   shell.Shell
}

func NewManager(conf *config.Config, sh shell.Shell) *Manager {
	return &Manager{
		conf: conf,
		sh:   sh,
	}
}

func (sm *Manager) Read() (*Settings, error) {
	settingsFileName := sm.conf.SettingsFileName

	settingsDirectory := filepath.Dir(settingsFileName)
	err := os.MkdirAll(settingsDirectory, os.ModePerm)
	if err != nil {
		return nil, err
	}

	exists, err := utils.Exists(settingsFileName)
	if err != nil {
		return nil, err
	}

	if !exists {
		newInstance := &Settings{}
		err := sm.Write(newInstance)
		if err != nil {
			return nil, err
		}

		return newInstance, nil
	}

	result := &Settings{}

	fileContent, err := os.ReadFile(settingsFileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileContent, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (sm *Manager) Write(settings *Settings) error {
	settingsFileName := sm.conf.SettingsFileName

	// make sure all data is sorted alphabetically
	slices.SortStableFunc(
		settings.Repos, func(a Repo, b Repo) int {
			return strings.Compare(a.Name, b.Name)
		},
	)
	slices.SortStableFunc(
		settings.Groups, func(a Group, b Group) int {
			return strings.Compare(a.Name, b.Name)
		},
	)
	for _, group := range settings.Groups {
		slices.Sort(group.Repos)
	}

	fileContent, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(settingsFileName), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(settingsFileName, fileContent, os.ModePerm)
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
	var fileModel *exportModel
	fileExists, err := utils.Exists(exportFileName)
	if err != nil {
		return nil, err
	}

	if fileExists {
		fileModel, err = readExistingModel(exportFileName)
		if err != nil {
			return nil, err
		}
	} else {
		fileModel = &exportModel{
			Version: currentVersion,
			Data:    modelDataV1{map[string]modelDataV1Repo{}},
		}
	}

	err = os.WriteFile(exportFileName, jsonBytes, os.ModePerm)
	if err != nil {
		return nil, err
	}

	_, err = sm.sh.RunCommand(repoDir, "git", "add", ".")
	if err != nil {
		return nil, err
	}

	statusOutput, err := sm.sh.RunCommand(repoDir, "git", "status")
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

	_, err = sm.sh.RunCommand(repoDir, "git", "commit", "-m", "Export bulker repositories")
	if err != nil {
		return nil, err
	}

	output, err := sm.sh.RunCommand(repoDir, "git", "push")
	if err != nil {
		logrus.Errorf("push failed: %v, %v", output, err)
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
	repoDirectory, err := os.MkdirTemp("", "bulker_remote_repo_*")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	cleanupFunc = func() {
		logrus.WithField("directory", repoDirectory).Debug("temporary directory deleted")
		err := os.RemoveAll(repoDirectory)
		if err != nil {
			logrus.Warnf("can't remove temp directory %v: %v", repoDirectory, err)
		}
	}
	logrus.WithField("repo", remoteRepoUrl).WithField("directory", repoDirectory).Debug("temporary directory created")

	output, err := sm.sh.RunCommand(repoDirectory, "git", "clone", remoteRepoUrl, ".")
	if err != nil {
		logrus.WithField("repo", remoteRepoUrl).WithField("output", output).Debug("clone failed")
		return "", cleanupFunc, fmt.Errorf("failed to clone repository: %w", err)
	}
	return repoDirectory, cleanupFunc, nil
}

func prepareResult(previousModel *exportModel, newModel *exportModel) map[string]ExportImportStatus {
	result := map[string]ExportImportStatus{}
	for repoName := range newModel.Data.Repos {
		if _, exists := previousModel.Data.Repos[repoName]; exists {
			result[repoName] = ExportImportStatusUpToDate
		} else {
			result[repoName] = ExportImportStatusAdded
		}
	}
	for repoName := range previousModel.Data.Repos {
		if _, exists := newModel.Data.Repos[repoName]; !exists {
			result[repoName] = ExportImportStatusRemoved
		}
	}
	return result
}

func readExistingModel(exportFileName string) (*exportModel, error) {
	fileBytes, err := os.ReadFile(exportFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file %s not found", exportFileName)
		}

		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var readVersionModel map[string]any
	err = yaml.Unmarshal(fileBytes, &readVersionModel)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall file content: %w", err)
	}

	version := readVersionModel["version"].(int)

	data := readVersionModel["data"]
	dataBytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall data: %w", err)
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

	return &exportModel{currentVersion, data}
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
