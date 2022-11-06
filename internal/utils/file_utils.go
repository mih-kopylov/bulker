package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func AbsPathify(path string) string {
	path = filepath.FromSlash(path)

	if path == "$HOME" || strings.HasPrefix(path, "$HOME"+string(os.PathSeparator)) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logrus.Fatalf("can't get home directory: %v", err)
		}
		path = homeDir + path[5:]
	}

	path = os.ExpandEnv(path)

	if !filepath.IsAbs(path) {
		path, err := filepath.Abs(path)
		if err != nil {
			logrus.WithField("path", path).Fatalf("can't make path absolute: %v", err)
		}
	}
	return filepath.Clean(path)
}

// Exists returns whether a file or directory exists on file system
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
