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
	"sync"
)

type Runner struct {
	fs      afero.Fs
	manager *settings.Manager
	config  *config.Config
	filter  *Filter
}

func NewRunner(fs afero.Fs, config *config.Config, filter *Filter) *Runner {
	return &Runner{
		fs:      fs,
		manager: settings.NewManager(fs, config),
		config:  config,
		filter:  filter,
	}
}

type RunContext struct {
	FS      afero.Fs
	Manager *settings.Manager
	Config  *config.Config
	Repo    *model.Repo
}

type ProcessResult struct {
	Result interface{}
	Error  error
}

type RepoHandler func(ctx context.Context, runContext *RunContext) (interface{}, error)

func (r Runner) Run(handler RepoHandler) error {
	sets, err := r.manager.Read()
	if err != nil {
		return err
	}

	ctx := context.Background()

	allReposResult := map[string]ProcessResult{}
	wg := sync.WaitGroup{}
	for _, repo := range sets.Repos {
		if !r.filter.Matches(repo) {
			continue
		}
		runContext := &RunContext{
			FS:      r.fs,
			Manager: r.manager,
			Config:  r.config,
			Repo: &model.Repo{
				Name: repo.Name,
				Path: path.Join(r.config.ReposDirectory, repo.Name),
				Url:  repo.Url,
			},
		}
		if r.config.RunMode == config.Sequential {
			repoResult, err := handler(ctx, runContext)
			allReposResult[runContext.Repo.Name] = ProcessResult{
				Result: repoResult,
				Error:  err,
			}
		} else if r.config.RunMode == config.Parallel {
			wg.Add(1)
			go func() {
				defer wg.Done()
				repoResult, err := handler(ctx, runContext)
				allReposResult[runContext.Repo.Name] = ProcessResult{
					Result: repoResult,
					Error:  err,
				}
			}()
		} else {
			return fmt.Errorf("unknown run mode: %v", r.config.RunMode)
		}
	}
	wg.Wait()

	valueToLog := map[string]output.EntityInfo{}
	for repoName, processResult := range allReposResult {
		valueToLog[repoName] = output.EntityInfo{
			Result: processResult.Result,
			Error:  processResult.Error,
		}
	}

	logrus.WithField("count", len(allReposResult)).Debug("processed repos")

	err = output.Write("repo", valueToLog)
	if err != nil {
		return err
	}

	return nil
}
