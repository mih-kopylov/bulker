package runner

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path"
)

type Runner interface {
	Run(handler RepoHandler) error
}

func NewRunner(fs afero.Fs, conf *config.Config, filter *Filter) (Runner, error) {
	manager := settings.NewManager(fs, conf)
	if conf.RunMode == config.Sequential {
		return &SequentialRunner{
			fs:      fs,
			manager: manager,
			config:  conf,
			filter:  filter,
		}, nil
	} else if conf.RunMode == config.Parallel {
		return &ParallelRunner{
			fs:      fs,
			manager: manager,
			config:  conf,
			filter:  filter,
		}, nil
	}
	return nil, fmt.Errorf("unsupported run mode %v", conf.RunMode)
}

type RunContext struct {
	FS      afero.Fs
	Manager *settings.Manager
	Config  *config.Config
	Repo    *model.Repo
}

func newRunContext(fs afero.Fs, manager *settings.Manager, conf *config.Config, repo settings.Repo) *RunContext {
	return &RunContext{
		FS:      fs,
		Manager: manager,
		Config:  conf,
		Repo: &model.Repo{
			Name: repo.Name,
			Path: path.Join(conf.ReposDirectory, repo.Name),
			Url:  repo.Url,
		},
	}
}

func logOutput(result map[string]ProcessResult) error {
	logrus.WithField("count", len(result)).Debug("processed repos")

	valueToLog := map[string]output.EntityInfo{}
	for repoName, procResult := range result {
		if procResult.Result == nil && procResult.Error == nil {
			continue
		}
		valueToLog[repoName] = output.EntityInfo{
			Result: procResult.Result,
			Error:  procResult.Error,
		}
	}

	err := output.Write("repo", valueToLog)
	if err != nil {
		return err
	}

	return nil
}

type ProcessResult struct {
	Result interface{}
	Error  error
}

type RepoHandler func(ctx context.Context, runContext *RunContext) (interface{}, error)
