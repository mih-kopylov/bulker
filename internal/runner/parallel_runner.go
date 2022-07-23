package runner

import (
	"context"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/mih-kopylov/bulker/internal/settings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"sync"
)

type ParallelRunner struct {
	fs      afero.Fs
	manager *settings.Manager
	config  *config.Config
	filter  *Filter
}

func (r ParallelRunner) Run(handler RepoHandler) error {
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
	wg := sync.WaitGroup{}
	ch := make(chan repoProcessResult)
	logrus.WithField("mode", r.config.RunMode).Debug("processing repositories in parallel")
	processedRepoCount := 0
	for _, repo := range sets.Repos {
		if !r.filter.Matches(repo) {
			continue
		}
		processedRepoCount++
		runContext := newRunContext(r.fs, r.manager, r.config, repo)
		wg.Add(1)
		go func(ch chan repoProcessResult) {
			defer wg.Done()
			repoResult, err := handler(ctx, runContext)
			ch <- repoProcessResult{
				Name: runContext.Repo.Name,
				ProcessResult: ProcessResult{
					Result: repoResult,
					Error:  err,
				},
			}
		}(ch)
	}
	for i := 0; i < processedRepoCount; i++ {
		result := <-ch
		allReposResult[result.Name] = result.ProcessResult
	}
	wg.Wait()
	close(ch)

	err = logOutput(allReposResult)
	if err != nil {
		return err
	}

	return nil
}
