package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Debug            bool         `mapstructure:"debug"`
	SettingsFileName string       `mapstructure:"settingsFileName"`
	ReposDirectory   string       `mapstructure:"reposDirectory"`
	RunMode          RunMode      `mapstructure:"runMode"`
	MaxWorkers       int          `mapstructure:"maxWorkers"`
	NoProgress       bool         `mapstructure:"noProgress"`
	Output           OutputFormat `mapstructure:"output"`
}

// RunMode implements Value in spf13/pflag for custom flag type
type RunMode string

const (
	Parallel   RunMode = "par"
	Sequential RunMode = "seq"
)

func (rm *RunMode) String() string {
	return string(*rm)
}

func (rm *RunMode) Set(v string) error {
	switch v {
	case string(Parallel), string(Sequential):
		*rm = RunMode(v)
		return nil
	default:
		return fmt.Errorf("must be either '%s' or '%s'", Parallel, Sequential)
	}
}

func (rm *RunMode) Type() string {
	return "RunMode"
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

	err = os.WriteFile(fileName, confBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
