package scheduler

import (
	"context"
	"time"
)

type Runnable interface {
	Run(ctx context.Context)
}

type Runner struct {
	interval time.Duration
	jobs     []Runnable
}

func New(interval time.Duration, jobs ...Runnable) *Runner {
	return &Runner{
		interval: interval,
		jobs:     jobs,
	}
}

func (r *Runner) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, job := range r.jobs {
				job.Run(ctx)
			}
		case <-ctx.Done():
			return
		}
	}
}
