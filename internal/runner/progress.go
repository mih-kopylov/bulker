package runner

import (
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

type Progress interface {
	// Incr increase current progress with 1
	Incr()
	// IndicateTermination makes progress indicate that bulker received "SIGINT" and terminates gracefully
	IndicateTermination()
}

func NewProgress(conf *config.Config, maxCount int) Progress {
	return newProgressBarProgress(maxCount, !conf.NoProgress && !conf.Debug)
}

type ProgressBarProgress struct {
	bar *progressbar.ProgressBar
}

func (p *ProgressBarProgress) Incr() {
	err := p.bar.Add(1)
	if err != nil {
		logrus.Debug(err)
	}
}

func (p *ProgressBarProgress) IndicateTermination() {
	p.bar.Describe("terminating, please wait")
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
