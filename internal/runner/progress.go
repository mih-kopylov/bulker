package runner

import (
	"fmt"
	"github.com/mih-kopylov/bulker/internal/config"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"strings"
)

type Progress interface {
	// IncrProgress increase current progress with 1
	IncrProgress()
	// IncrErrors increments number of errors happened during the process
	IncrErrors()
	// IndicateTermination makes progress indicate that bulker received "SIGINT" and terminates gracefully
	IndicateTermination()
}

func NewProgress(conf *config.Config, maxCount int) Progress {
	return newProgressBarProgress(maxCount, !conf.NoProgress && !conf.Debug)
}

type ProgressBarProgress struct {
	bar         *progressbar.ProgressBar
	errorsCount int
	terminating bool
}

func (p *ProgressBarProgress) IncrProgress() {
	err := p.bar.Add(1)
	if err != nil {
		logrus.Debug(err)
	}
}

func (p *ProgressBarProgress) IncrErrors() {
	p.errorsCount += 1
	p.updateDescription()
}

func (p *ProgressBarProgress) IndicateTermination() {
	p.terminating = true
	p.updateDescription()
}

func (p *ProgressBarProgress) updateDescription() {
	var messages []string
	if p.errorsCount > 0 {
		messages = append(messages, fmt.Sprintf("errors: %v", p.errorsCount))
	}
	if p.terminating {
		messages = append(messages, "terminating, please wait")
	}
	p.bar.Describe(strings.Join(messages, " | "))

}

func newProgressBarProgress(maxCount int, visible bool) *ProgressBarProgress {
	return &ProgressBarProgress{
		bar: progressbar.NewOptions(
			maxCount,
			progressbar.OptionFullWidth(),
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionShowCount(),
			progressbar.OptionSetVisibility(visible),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionClearOnFinish(),
		),
		errorsCount: 0,
	}
}
