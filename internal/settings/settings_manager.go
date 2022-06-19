package settings

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"os"
)

type Manager struct {
	fs afero.Fs
}

func NewManager(fs afero.Fs) *Manager {
	return &Manager{
		fs: fs,
	}
}

func (sm *Manager) Read() (*Settings, error) {
	settingsFileName := config.ReadConfig().SettingsFileName

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
	settingsFileName := config.ReadConfig().SettingsFileName

	fileContent, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	err = afero.WriteFile(sm.fs, settingsFileName, fileContent, os.ModePerm)
	return err
}
