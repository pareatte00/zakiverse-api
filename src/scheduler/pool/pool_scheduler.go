package pool

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

	// 1. Find pools due for rotation
	pools, err := s.repository.PackPool.FindDueRotations(ctx, now)
	if err != nil {
		log.Printf("[pool-scheduler] error finding due rotations: %v", err)
		return
	}

	for _, pool := range pools {
		poolId := pool.ID.String()

		// 2. Deactivate all active packs in the pool
		err := s.repository.Pack.DeactivateByPool(ctx, poolId)
		if err != nil {
			log.Printf("[pool-scheduler] error deactivating packs for pool %s: %v", poolId, err)
			continue
		}

		// 3. Find all packs in pool ordered by last_pool_activated_at ASC NULLS FIRST
		packs, err := s.repository.Pack.FindByPool(ctx, poolId)
		if err != nil {
			log.Printf("[pool-scheduler] error finding packs for pool %s: %v", poolId, err)
			continue
		}

		if len(packs) == 0 {
			continue
		}

		// 4. Activate top active_count packs
		activateCount := int(pool.ActiveCount)
		if activateCount > len(packs) {
			activateCount = len(packs)
		}

		for i := 0; i < activateCount; i++ {
			packId := packs[i].ID.String()
			err := s.repository.Pack.ActivateWithPoolTimestamp(ctx, packId, now)
			if err != nil {
				log.Printf("[pool-scheduler] error activating pack %s: %v", packId, err)
				continue
			}
			log.Printf("[pool-scheduler] activated pack %s (%s) in pool %s", packId, packs[i].Name, pool.Name)
		}

		// 5. Update pool's last_rotated_at
		err = s.repository.PackPool.SetLastRotated(ctx, poolId, now)
		if err != nil {
			log.Printf("[pool-scheduler] error updating last_rotated_at for pool %s: %v", poolId, err)
		}
	}
}
