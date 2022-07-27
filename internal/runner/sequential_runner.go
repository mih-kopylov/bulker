package runner

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type SequentialRunner struct {
	fs      afero.Fs
	manager *settings.Manager
	config  *config.Config
	filter  *Filter
}

func (r SequentialRunner) Run(handler RepoHandler) error {
	sets, err := r.manager.Read()
	if err != nil {
		return err
	}

	ctx := context.Background()

	allReposResult := map[string]ProcessResult{}
	logrus.WithField("mode", r.config.RunMode).Debug("processing repositories")
	for _, repo := range sets.Repos {
		if !r.filter.Matches(repo) {
			continue
		}
		runContext := newRunContext(r.fs, r.manager, r.config, repo)
		logrus.WithField("repo", repo.Name).Debug("processing started")
		repoResult, err := handler(ctx, runContext)
		logrus.WithField("repo", repo.Name).Debug("processing completed")
		allReposResult[runContext.Repo.Name] = ProcessResult{
			Result: repoResult,
			Error:  err,
		}
	}

	err = logOutput(allReposResult)
	if err != nil {
		return err
	}

	return nil
}
