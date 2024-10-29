package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	config2 "github.com/go-git/go-git/v5/config"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func ExecuteCommand(root *cobra.Command, args string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)

	root.SetArgs(strings.Split(args, " "))

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func PrepareBulker(t *testing.T, sh shell.Shell, repos []settings.Repo) {
	PrepareBulkerWithGroups(t, sh, repos, nil)
}

func PrepareBulkerWithGroups(t *testing.T, sh shell.Shell, repos []settings.Repo, groups []settings.Group) {
	testDirectory := t.TempDir()
	viper.Set("reposDirectory", testDirectory)
	viper.Set("settingsFileName", filepath.Join(testDirectory, "test_bulker_settings.yaml"))
	viper.Set("runMode", "seq")
	viper.Set("noProgress", "true")
	viper.Set("output", "json")

	conf := config.ReadConfig()
	manager := settings.NewManager(conf, sh)
	s := &settings.Settings{
		Repos:  repos,
		Groups: groups,
	}
	err := manager.Write(s)
	assert.NoError(t, err)
}

func ShellCommandToString(command string, arguments []string) string {
	buffer := bytes.Buffer{}
	buffer.WriteString(command)
	buffer.WriteString(" ")
	buffer.WriteString(strings.Join(arguments, " "))
	return buffer.String()
}

func ToJsonString(value any) string {
	marshalled, err := json.Marshal(value)
	if err != nil {
		logrus.Fatal(err)
	}
	return string(marshalled)
}

func Path(parts ...string) string {
	allParts := slices.Insert(parts, 0, viper.GetString("reposDirectory"))
	return filepath.Join(allParts...)
}

func CreateBareGitRepository(repoName string) (string, error) {
	repoPath := Path(repoName)
	repo, err := git.PlainInit(repoPath, true)
	if err != nil {
		return "", err
	}

	conf, err := repo.Config()
	if err != nil {
		return "", err
	}

	conf.User.Name = "test"
	conf.User.Email = "test@example.com"
	err = repo.SetConfig(conf)
	if err != nil {
		return "", err
	}

	//print configs
	localConfig, _ := repo.ConfigScoped(config2.LocalScope)
	println("local", localConfig.User.Name, localConfig.User.Email)
	globalConfig, _ := repo.ConfigScoped(config2.GlobalScope)
	println("global", globalConfig.User.Name, globalConfig.User.Email)
	systemConfig, _ := repo.ConfigScoped(config2.SystemScope)
	println("system", systemConfig.User.Name, systemConfig.User.Email)

	return repoPath, nil
}

func GetFileContent(repoPath string, fileName string) (string, error) {
	repository, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open git repo: repoPath=%v, err=%w", repoPath, err)
	}

	headReference, err := repository.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get repo HEAD: repoPath=%v, err=%w", repoPath, err)
	}

	headHash := headReference.Hash()

	commit, err := repository.CommitObject(headHash)
	if err != nil {
		return "", fmt.Errorf(
			"failed to find HEAD commit: repoPath=%v, commit=%v, err=%w", repoPath,
			headHash.String(), err,
		)
	}

	file, err := commit.File(fileName)
	if err != nil {
		return "", fmt.Errorf(
			"failed to find file in commit: repoPath=%v, commit=%v, file=%v, err=%w", repoPath,
			headHash.String(), fileName, err,
		)
	}

	fileContent, err := file.Contents()
	if err != nil {
		return "", fmt.Errorf(
			"failed to get file content: repoPath=%v, commit=%v, file=%v, err=%w", repoPath,
			headHash.String(), fileName, err,
		)
	}

	return fileContent, nil
}
