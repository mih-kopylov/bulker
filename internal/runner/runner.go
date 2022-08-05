package runner

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"path/filepath"
	"time"
)

type Runner interface {
	Run(ctx context.Context, repos []settings.Repo, handler RepoHandler) (map[string]ProcessResult, error)
}

func NewRunner(fs afero.Fs, conf *config.Config, filter *Filter, progress Progress, args []string) (Runner, error) {
	manager := settings.NewManager(fs, conf)
	if conf.RunMode == config.Sequential {
		return &SequentialRunner{
			fs:       fs,
			manager:  manager,
			config:   conf,
			filter:   filter,
			progress: progress,
			args:     args,
		}, nil
	} else if conf.RunMode == config.Parallel {
		return &ParallelRunner{
			fs:       fs,
			manager:  manager,
			config:   conf,
			filter:   filter,
			progress: progress,
			args:     args,
		}, nil
	}
	return nil, fmt.Errorf("unsupported run mode %v", conf.RunMode)
}

func NewDefaultRunner(filter *Filter, handler RepoHandler) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fs := utils.GetConfiguredFS()
		conf := config.ReadConfig()
		manager := settings.NewManager(fs, conf)

		sets, err := manager.Read()
		if err != nil {
			return err
		}

		repos := filter.FilterMatchingRepos(sets.Repos, sets.Groups)
		progress := NewProgress(conf, len(repos))
		newRunner, err := NewRunner(fs, conf, filter, progress, args)
		if err != nil {
			return err
		}

		go func() {
			for {
				select {
				case <-cmd.Context().Done():
					progress.IndicateTermination()
					return
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()

		allReposResult, err := newRunner.Run(cmd.Context(), repos, handler)
		if err != nil {
			return err
		}

		err = logOutput(allReposResult)
		if err != nil {
			return err
		}

		return nil
	}
}

type RunContext struct {
	Fs      afero.Fs
	Manager *settings.Manager
	Config  *config.Config
	Repo    *model.Repo
	Args    []string
}

func newRunContext(
	fs afero.Fs, manager *settings.Manager, conf *config.Config, args []string, repo settings.Repo,
) *RunContext {
	return &RunContext{
		Fs:      fs,
		Manager: manager,
		Config:  conf,
		Repo: &model.Repo{
			Name: repo.Name,
			Path: filepath.Join(conf.ReposDirectory, repo.Name),
			Url:  repo.Url,
		},
		Args: args,
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
