package runner

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/schollz/progressbar/v3"
)

type Progress interface {
	Incr()
}

func NewProgress(conf *config.Config, maxCount int) Progress {
	return newProgressBarProgress(maxCount, !conf.NoProgress && !conf.Debug)
}

type ProgressBarProgress struct {
	bar *progressbar.ProgressBar
}

func (p *ProgressBarProgress) Incr() {
	_ = p.bar.Add(1)
}

func newProgressBarProgress(maxCount int, visible bool) *ProgressBarProgress {
	return &ProgressBarProgress{
		progressbar.NewOptions(
			maxCount,
			progressbar.OptionFullWidth(),
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionShowCount(),
			progressbar.OptionSetVisibility(visible),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionClearOnFinish(),
		),
	}
}
