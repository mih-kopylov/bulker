package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func MarkFlagRequiredOrFail(flagSet *pflag.FlagSet, name string) {
	err := cobra.MarkFlagRequired(flagSet, name)
	if err != nil {
		logrus.Fatalf("failed to mark flag required: %v", err)
	}
}

func BindFlag(flag *pflag.Flag, viperFlagName string) {
	err := viper.BindPFlag(viperFlagName, flag)
	if err != nil {
		logrus.Fatalf("can't bind flag: %v", err)
	}

}
