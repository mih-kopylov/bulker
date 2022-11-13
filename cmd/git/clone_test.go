package git

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/mih-kopylov/bulker/internal/tests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestClone(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := shell.MockShell(
		func(repoName string, command string, arguments []string) (string, error) {
			if tests.ShellCommandEquals(command, arguments, "git status") {
				return "OK", nil
			}
			if tests.ShellCommandEquals(command, arguments, "git clone https://example.com .") {
				return "OK", nil
			}

			return "", fmt.Errorf("not mocked")
		},
	)
	tests.PrepareBulker(t, sh, repos)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n", "repo")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.Equal(t, output, "repo: result=Cloned\n")
}

func TestClone_EmptyDirectoryExists(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := shell.MockShell(
		func(repoName string, command string, arguments []string) (string, error) {
			if tests.ShellCommandEquals(command, arguments, "git status") {
				return "OK", nil
			}
			if tests.ShellCommandEquals(command, arguments, "git clone https://example.com .") {
				return "OK", nil
			}

			return "", fmt.Errorf("not mocked")
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n", "repo")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.Equal(t, output, "repo: result=Cloned\n")
}

func TestClone_NotEmptyDirectoryExists(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := shell.MockShell(
		func(repoName string, command string, arguments []string) (string, error) {
			if tests.ShellCommandEquals(command, arguments, "git status") {
				return "err", fmt.Errorf("not a repository")
			}

			return "", fmt.Errorf("not mocked")
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(viper.GetString("reposDirectory"), "repo", "file"), []byte("data"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n", "repo")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.Equal(t, output, "repo: failed to clone: failed to get git status: err, not a repository\n")
}

func TestClone_Recreate(t *testing.T) {
	repos := []settings.Repo{
		{
			Name: "repo",
			Url:  "https://example.com",
		},
	}
	sh := shell.MockShell(
		func(repoName string, command string, arguments []string) (string, error) {
			if tests.ShellCommandEquals(command, arguments, "git status") {
				return "On branch custom\nYour branch is up to date with 'origin/custom'.", nil
			}
			if tests.ShellCommandEquals(command, arguments, "git clone https://example.com .") {
				return "OK", nil
			}

			return "", fmt.Errorf("not mocked")
		},
	)
	tests.PrepareBulker(t, sh, repos)
	err := os.Mkdir(filepath.Join(viper.GetString("reposDirectory"), "repo"), os.ModePerm)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(viper.GetString("reposDirectory"), "repo", "file"), []byte("data"), os.ModePerm)
	assert.NoError(t, err)

	command := CreateCloneCommand(sh)
	c, output, err := tests.ExecuteCommand(command, "-n", "repo", "--recreate")
	assert.NoError(t, err)
	assert.Equal(t, "clone", c.Name())
	assert.Equal(t, output, "repo: result=Re-cloned\n")
}
