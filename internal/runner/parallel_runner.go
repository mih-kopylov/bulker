package runner

import (
	"context"
	"github.com/alitto/pond"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type ParallelRunner struct {
	fs      afero.Fs
	manager *settings.Manager
	config  *config.Config
	filter  *Filter
}

func (r *ParallelRunner) Run(handler RepoHandler) error {
	type repoProcessResult struct {
		Name string
		ProcessResult
	}

	sets, err := r.manager.Read()
	if err != nil {
		return err
	}

	ctx := context.Background()

	allReposResult := map[string]ProcessResult{}
	pool := pond.New(r.config.MaxWorkers, 1000, pond.Context(ctx))
	defer pool.StopAndWait()
	ch := make(chan repoProcessResult)
	logrus.WithField("mode", r.config.RunMode).Debug("processing repositories")
	repos := r.filter.FilterMatchingRepos(sets.Repos, sets.Groups)
	progress := NewProgress(r.config, len(repos))
	for _, repo := range repos {
		runContext := newRunContext(r.fs, r.manager, r.config, repo)
		pool.Submit(
			func() {
				logrus.WithField("repo", runContext.Repo.Name).Debug("processing started")
				repoResult, err := handler(ctx, runContext)
				progress.Incr()
				logrus.WithField("repo", runContext.Repo.Name).Debug("processing completed")
				ch <- repoProcessResult{
					Name: runContext.Repo.Name,
					ProcessResult: ProcessResult{
						Result: repoResult,
						Error:  err,
					},
				}
			},
		)
	}
	for i := 0; i < len(repos); i++ {
		result := <-ch
		allReposResult[result.Name] = result.ProcessResult
	}
	close(ch)

	err = logOutput(allReposResult)
	if err != nil {
		return err
	}

	return nil
}
