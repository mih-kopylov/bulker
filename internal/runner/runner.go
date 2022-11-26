package runner

import (
	"context"
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/fileops"
	"github.com/mih-kopylov/bulker/internal/model"
	"github.com/mih-kopylov/bulker/internal/output"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/mih-kopylov/bulker/internal/shell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"io"
	"path/filepath"
	"time"
)

type Runner interface {
	Run(ctx context.Context, repos []settings.Repo, handler RepoHandler) (map[string]ProcessResult, error)
}

func NewRunner(conf *config.Config, sh shell.Shell, filter *Filter, progress Progress, args []string) (Runner, error) {
	manager := settings.NewManager(conf, sh)
	if conf.RunMode == config.Sequential {
		return &SequentialRunner{
			manager:  manager,
			config:   conf,
			filter:   filter,
			progress: progress,
			args:     args,
		}, nil
	} else if conf.RunMode == config.Parallel {
		return &ParallelRunner{
			manager:  manager,
			config:   conf,
			filter:   filter,
			progress: progress,
			args:     args,
		}, nil
	}
	return nil, fmt.Errorf("unsupported run mode %v", conf.RunMode)
}

func NewCommandRunner(filter *Filter, sh shell.Shell, handler RepoHandler) func(
	cmd *cobra.Command,
	args []string,
) error {
	return func(cmd *cobra.Command, args []string) error {
		conf := config.ReadConfig()
		manager := settings.NewManager(conf, sh)

		sets, err := manager.Read()
		if err != nil {
			return err
		}

		repos := filter.FilterMatchingRepos(sets.Repos, sets.Groups)
		progress := NewProgress(conf, len(repos))
		newRunner, err := NewRunner(conf, sh, filter, progress, args)
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

		err = savePreviousGroup(manager, maps.Keys(allReposResult))
		if err != nil {
			return err
		}

		err = logOutput(cmd.OutOrStdout(), allReposResult)
		if err != nil {
			return err
		}

		return nil
	}
}

func NewCommandRunnerForExistingRepos(filter *Filter, sh shell.Shell, handler RepoHandler) func(
	cmd *cobra.Command,
	args []string,
) error {
	handlerWithRepoExistenceVerifier := func(ctx context.Context, runContext *RunContext) (interface{}, error) {
		err := fileops.CheckRepoExists(runContext.Repo)
		if err != nil {
			return nil, err
		}

		return handler(ctx, runContext)
	}

	return NewCommandRunner(filter, sh, handlerWithRepoExistenceVerifier)
}

type RunContext struct {
	Manager *settings.Manager
	Config  *config.Config
	Repo    *model.Repo
	Args    []string
}

func newRunContext(manager *settings.Manager, conf *config.Config, args []string, repo settings.Repo) *RunContext {
	return &RunContext{
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

// savePreviousGroup saves a group with a constant name `previous`. If such one exists, it gets recreated
func savePreviousGroup(manager *settings.Manager, repos []string) error {
	sets, err := manager.Read()
	if err != nil {
		return err
	}

	if sets.GroupExists(settings.PreviousGroupName) {
		err := sets.RemoveGroup(settings.PreviousGroupName)
		if err != nil {
			return err
		}
	}

	group, err := sets.AddGroup(settings.PreviousGroupName)
	if err != nil {
		return err
	}

	for _, repoName := range repos {
		err := sets.AddRepoToGroup(group, repoName)
		if err != nil {
			return err
		}
	}

	err = manager.Write(sets)
	if err != nil {
		return err
	}

	return nil
}

func logOutput(writer io.Writer, result map[string]ProcessResult) error {
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

	err := output.Write(writer, "repo", valueToLog)
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
