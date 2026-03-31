package scheduler

import (
	"context"
	"time"

	"github.com/zakiverse/zakiverse-api/src/repository"
	"github.com/zakiverse/zakiverse-api/src/scheduler/pack"
)

const (
	t1m = 1 * time.Minute
)

func Start(ctx context.Context, repository *repository.Repository) {
	go New(t1m, pack.NewScheduler(repository)).Start(ctx)
}
