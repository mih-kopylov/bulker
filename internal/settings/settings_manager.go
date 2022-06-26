package settings

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"os"
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
