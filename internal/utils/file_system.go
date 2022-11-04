package utils

import (
	"github.com/spf13/afero"
)

var fs afero.Fs

// GetConfiguredFS Gets the configured file system for the application. Once it's configured on start,
// should be used for all file system manipulations
func GetConfiguredFS() afero.Fs {
	return fs
}

// ConfigureFS Configures file system to be used in the application. Looks like a singleton, but not a real one
func ConfigureFS(newFs afero.Fs) {
	fs = newFs
}
