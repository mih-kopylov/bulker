package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Debug            bool         `mapstructure:"debug"`
	SettingsFileName string       `mapstructure:"settings"`
	ReposDirectory   string       `mapstructure:"reposDirectory"`
	RunMode          RunMode      `mapstructure:"runMode"`
	MaxWorkers       int          `mapstructure:"maxWorkers"`
	Output           OutputFormat `mapstructure:"output"`
}

// implements Value in spf13/pflag for custom flag type
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

type OutputFormat string

const (
	JsonOutputFormat OutputFormat = "json"
	LineOutputFormat OutputFormat = "line"
	LogOutputFormat  OutputFormat = "log"
)

func ReadConfig() *Config {
	config := &Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		logrus.Fatalf("can't read config: %v", err)
	}

	return config
}
