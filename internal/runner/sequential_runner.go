package runner

import (
	"context"
	"errors"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type SequentialRunner struct {
	fs       afero.Fs
	manager  *settings.Manager
	config   *config.Config
	filter   *Filter
	progress Progress
	args     []string
}

func (r *SequentialRunner) Run(
	ctx context.Context, repos []settings.Repo, handler RepoHandler,
) (map[string]ProcessResult, error) {
	allReposResult := map[string]ProcessResult{}
	logrus.WithField("mode", r.config.RunMode).Debug("processing repositories")
	for _, repo := range repos {
		runContext := newRunContext(r.fs, r.manager, r.config, r.args, repo)
		select {
		case <-ctx.Done():
			logrus.WithField("repo", runContext.Repo.Name).Debug("processing skipped")
			r.progress.IncrProgress()
			allReposResult[runContext.Repo.Name] = ProcessResult{
				Result: nil,
				Error:  errors.New("skipped"),
			}
		default:
			logrus.WithField("repo", repo.Name).Debug("processing started")
			repoResult, err := handler(ctx, runContext)
			logrus.WithField("repo", repo.Name).Debug("processing completed")
			r.progress.IncrProgress()
			if err != nil {
				r.progress.IncrErrors()
			}
			allReposResult[runContext.Repo.Name] = ProcessResult{
				Result: repoResult,
				Error:  err,
			}
		}
	}

	return allReposResult, nil
}
