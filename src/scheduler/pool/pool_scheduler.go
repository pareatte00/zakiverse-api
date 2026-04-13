package pool

import (
	"context"
	"log"
	"time"

	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	"github.com/zakiverse/zakiverse-api/src/repository"
)

type Scheduler struct {
	repository *repository.Repository
}

func NewScheduler(repository *repository.Repository) *Scheduler {
	return &Scheduler{repository: repository}
}

func (s *Scheduler) Run(ctx context.Context) {
	now := time.Now().UTC()

	s.activatePendingPools(ctx, now)
	s.deactivateExpiredPools(ctx, now)
	s.rotateDuePools(ctx, now)
}

func (s *Scheduler) activatePendingPools(ctx context.Context, now time.Time) {
	pools, err := s.repository.PackPool.FindPendingOpen(ctx, now)
	if err != nil {
		log.Printf("[pool-scheduler] error finding pending pools: %v", err)
		return
	}

	for _, pool := range pools {
		id := pool.ID.String()
		if err := s.repository.PackPool.SetActive(ctx, id, true); err != nil {
			log.Printf("[pool-scheduler] error activating pool %s: %v", id, err)
			continue
		}
		log.Printf("[pool-scheduler] activated pool %s (%s)", id, pool.Name)
	}
}

func (s *Scheduler) deactivateExpiredPools(ctx context.Context, now time.Time) {
	pools, err := s.repository.PackPool.FindExpired(ctx, now)
	if err != nil {
		log.Printf("[pool-scheduler] error finding expired pools: %v", err)
		return
	}

	for _, pool := range pools {
		id := pool.ID.String()
		if err := s.repository.PackPool.SetActive(ctx, id, false); err != nil {
			log.Printf("[pool-scheduler] error deactivating pool %s: %v", id, err)
			continue
		}
		log.Printf("[pool-scheduler] deactivated pool %s (%s)", id, pool.Name)
	}
}

func (s *Scheduler) rotateDuePools(ctx context.Context, now time.Time) {
	pools, err := s.repository.PackPool.FindDueRotations(ctx, now)
	if err != nil {
		log.Printf("[pool-scheduler] error finding due rotations: %v", err)
		return
	}

	for _, pool := range pools {
		poolId := pool.ID.String()

		// Get packs ordered for rotation based on mode
		var packs []model.Pack
		if pool.RotationOrderMode == model.RotationOrderMode_Manual {
			packs, err = s.repository.Pack.FindNextRotationByPool(ctx, poolId, pool.ActiveCount)
		} else {
			packs, err = s.repository.Pack.FindByPool(ctx, poolId)
		}
		if err != nil {
			log.Printf("[pool-scheduler] error finding packs for pool %s: %v", poolId, err)
			continue
		}

		if len(packs) == 0 {
			continue
		}

		// Activate top active_count packs
		activateCount := min(int(pool.ActiveCount), len(packs))

		for i := range activateCount {
			packId := packs[i].ID.String()
			if err := s.repository.Pack.SetPoolActivatedAt(ctx, packId, now); err != nil {
				log.Printf("[pool-scheduler] error activating pack %s: %v", packId, err)
				continue
			}
			log.Printf("[pool-scheduler] rotated pack %s (%s) in pool %s", packId, packs[i].Name, pool.Name)
		}

		// Compute next rotation and update pool
		nextRotation := computeNextRotationAt(pool, now)
		if err := s.repository.PackPool.SetLastRotated(ctx, poolId, now, nextRotation); err != nil {
			log.Printf("[pool-scheduler] error updating last_rotated_at for pool %s: %v", poolId, err)
		}
	}
}

// computeNextRotationAt advances the existing next_rotation_at (already in UTC) by the pool's interval.
// For weekly: adds rotation_interval * 7 days.
// For monthly: advances to same day next month (clamped to last day), preserving the UTC hour.
func computeNextRotationAt(pool model.PackPool, now time.Time) *time.Time {
	if pool.NextRotationAt == nil {
		return nil
	}
	prev := *pool.NextRotationAt

	switch pool.RotationType {
	case model.RotationType_Weekly:
		next := prev.AddDate(0, 0, int(pool.RotationInterval)*7)
		return &next

	case model.RotationType_Monthly:
		day := 1
		if pool.RotationDay != nil {
			day = int(*pool.RotationDay)
		}
		next := clampedMonthDay(prev.Year(), prev.Month()+1, day, prev.Hour())
		return &next
	}

	return nil
}

// clampedMonthDay creates a time for the given year/month/day, clamping the day
// to the last day of the month if it exceeds it.
func clampedMonthDay(year int, month time.Month, day int, hour int) time.Time {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(t.Year(), t.Month(), day, hour, 0, 0, 0, time.UTC)
}
