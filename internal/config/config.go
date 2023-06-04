package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Debug            bool         `mapstructure:"debug"`
	SettingsFileName string       `mapstructure:"settingsFileName"`
	ReposDirectory   string       `mapstructure:"reposDirectory"`
	RunMode          RunMode      `mapstructure:"runMode"`
	MaxWorkers       int          `mapstructure:"maxWorkers"`
	NoProgress       bool         `mapstructure:"noProgress"`
	Output           OutputFormat `mapstructure:"output"`
	GitMode          GitMode      `mapstructure:"gitMode"`
}

func ReadConfig() *Config {
	config := &Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		logrus.Fatalf("can't read config: %v", err)
	}

	return config
}

func WriteConfig(conf *Config, fileName string) error {
	var confMap map[string]any
	err := mapstructure.Decode(conf, &confMap)
	if err != nil {
		return err
	}

	confBytes, err := yaml.Marshal(confMap)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, confBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
