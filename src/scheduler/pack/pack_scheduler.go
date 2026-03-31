package pack

import (
	"context"
	"log"
	"time"

	"github.com/zakiverse/zakiverse-api/src/repository"
)

type Scheduler struct {
	repository *repository.Repository
}

func NewScheduler(repository *repository.Repository) *Scheduler {
	return &Scheduler{repository: repository}
}

func (s *Scheduler) Run(ctx context.Context) {
	now := time.Now()

	// Activate packs where open_at has passed
	pendingPacks, err := s.repository.Pack.FindPendingOpen(ctx, now)
	if err != nil {
		log.Printf("[pack-scheduler] error finding pending packs: %v", err)
		return
	}

	for _, p := range pendingPacks {
		err := s.repository.Pack.SetActive(ctx, p.ID.String(), true)
		if err != nil {
			log.Printf("[pack-scheduler] error activating pack %s: %v", p.ID, err)
			continue
		}
		log.Printf("[pack-scheduler] activated pack %s (%s)", p.ID, p.Name)
	}

	// Deactivate packs where close_at has passed
	expiredPacks, err := s.repository.Pack.FindExpired(ctx, now)
	if err != nil {
		log.Printf("[pack-scheduler] error finding expired packs: %v", err)
		return
	}

	for _, p := range expiredPacks {
		err := s.repository.Pack.SetActive(ctx, p.ID.String(), false)
		if err != nil {
			log.Printf("[pack-scheduler] error deactivating pack %s: %v", p.ID, err)
			continue
		}
		log.Printf("[pack-scheduler] deactivated pack %s (%s)", p.ID, p.Name)
	}
}
