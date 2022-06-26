package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Debug            bool     `mapstructure:"debug"`
	SettingsFileName string   `mapstructure:"settings"`
	ReposDirectory   string   `mapstructure:"reposDirectory"`
	RunMode          RunModes `mapstructure:"runMode"`
}

type RunModes string

const (
	Parallel   RunModes = "par"
	Sequential RunModes = "seq"
)

func ReadConfig() *Config {
	config := &Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		logrus.Fatalf("can't read config: %v", err)
	}

	return config
}
