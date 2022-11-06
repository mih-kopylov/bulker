package runner

import (
	"context"
	"errors"
	"github.com/alitto/pond"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/sirupsen/logrus"
)

type ParallelRunner struct {
	manager  *settings.Manager
	config   *config.Config
	filter   *Filter
	progress Progress
	args     []string
}

func (r *ParallelRunner) Run(
	ctx context.Context, repos []settings.Repo, handler RepoHandler,
) (map[string]ProcessResult, error) {
	type repoProcessResult struct {
		Name string
		ProcessResult
	}

	allReposResult := map[string]ProcessResult{}
	// the pool doesn't use the parent context in order to complete all the tasks even the context is done
	// each task is notified about the context is done on its own and passes skipped status to the output channel
	pool := pond.New(r.config.MaxWorkers, len(repos))
	defer pool.StopAndWait()
	ch := make(chan repoProcessResult)
	logrus.WithField("mode", r.config.RunMode).Debug("processing repositories")
	for _, repo := range repos {
		runContext := newRunContext(r.manager, r.config, r.args, repo)
		pool.Submit(
			func() {
				select {
				case <-ctx.Done():
					logrus.WithField("repo", runContext.Repo.Name).Debug("processing skipped")
					r.progress.IncrProgress()
					ch <- repoProcessResult{
						Name: runContext.Repo.Name,
						ProcessResult: ProcessResult{
							Result: nil,
							Error:  errors.New("skipped"),
						},
					}
				default:
					logrus.WithField("repo", runContext.Repo.Name).Debug("processing started")
					repoResult, err := handler(ctx, runContext)
					logrus.WithField("repo", runContext.Repo.Name).Debug("processing completed")
					r.progress.IncrProgress()
					if err != nil {
						r.progress.IncrErrors()
					}
					ch <- repoProcessResult{
						Name: runContext.Repo.Name,
						ProcessResult: ProcessResult{
							Result: repoResult,
							Error:  err,
						},
					}
				}
			},
		)
	}
	for i := 0; i < len(repos); i++ {
		result := <-ch
		allReposResult[result.Name] = result.ProcessResult
	}
	close(ch)

	return allReposResult, nil
}
