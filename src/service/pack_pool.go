package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	packRepo "github.com/zakiverse/zakiverse-api/src/repository/pack"
	poolRepo "github.com/zakiverse/zakiverse-api/src/repository/pack_pool"
	"github.com/zakiverse/zakiverse-api/util/pagination"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type PackPoolService struct {
	service *Service
}

type PackPoolPayload struct {
	ID                uuid.UUID          `json:"id"`
	Name              string             `json:"name"`
	Description       *string            `json:"description"`
	Image             *string            `json:"image"`
	BannerType        string             `json:"banner_type"`
	SortOrder         int32              `json:"sort_order"`
	IsActive          bool               `json:"is_active"`
	OpenAt            *time.Time         `json:"open_at"`
	CloseAt           *time.Time         `json:"close_at"`
	ActiveCount       int32              `json:"active_count"`
	RotationType      string             `json:"rotation_type"`
	RotationDay       *int32             `json:"rotation_day"`
	RotationInterval  int32              `json:"rotation_interval"`
	RotationHour      int32              `json:"rotation_hour"`
	RotationOrderMode string             `json:"rotation_order_mode"`
	NextRotationAt    *time.Time         `json:"next_rotation_at"`
	LastRotatedAt     *time.Time         `json:"last_rotated_at"`
	PreviewDays       int32              `json:"preview_days"`
	IsPreview         bool               `json:"is_preview"`
	HasPreview        bool               `json:"has_preview"`
	Packs             []PackPoolPackItem `json:"packs,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

type PackPoolPackItem struct {
	ID            uuid.UUID  `json:"id"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	Image         string     `json:"image"`
	NameImage     *string    `json:"name_image"`
	CardsPerPull  int32      `json:"cards_per_pull"`
	SortOrder     int32      `json:"sort_order"`
	Config        PackConfig `json:"config"`
	PoolId        *uuid.UUID `json:"pool_id"`
	RotationOrder *int32     `json:"rotation_order"`
	TotalCards    int64      `json:"total_cards"`
}

func toPackPoolPayload(pool model.PackPool) PackPoolPayload {
	return PackPoolPayload{
		ID:                pool.ID,
		Name:              pool.Name,
		Description:       pool.Description,
		Image:             pool.Image,
		BannerType:        string(pool.BannerType),
		SortOrder:         pool.SortOrder,
		IsActive:          pool.IsActive,
		OpenAt:            pool.OpenAt,
		CloseAt:           pool.CloseAt,
		ActiveCount:       pool.ActiveCount,
		RotationType:      string(pool.RotationType),
		RotationDay:       pool.RotationDay,
		RotationInterval:  pool.RotationInterval,
		RotationHour:      pool.RotationHour,
		RotationOrderMode: string(pool.RotationOrderMode),
		NextRotationAt:    pool.NextRotationAt,
		LastRotatedAt:     pool.LastRotatedAt,
		PreviewDays:       pool.PreviewDays,
		CreatedAt:         pool.CreatedAt,
		UpdatedAt:         pool.UpdatedAt,
	}
}

func toPackPoolPackItems(packs []packRepo.PackWithCardCount) []PackPoolPackItem {
	items := make([]PackPoolPackItem, len(packs))
	for i, p := range packs {
		items[i] = PackPoolPackItem{
			ID:            p.ID,
			Code:          p.Code,
			Name:          p.Name,
			Description:   p.Description,
			Image:         p.Image,
			NameImage:     p.NameImage,
			CardsPerPull:  p.CardsPerPull,
			SortOrder:     p.SortOrder,
			Config:        unmarshalPackConfig(p.Config),
			PoolId:        p.PoolID,
			RotationOrder: p.RotationOrder,
			TotalCards:    p.TotalCards,
		}
	}
	return items
}

// computeNextRotationAt calculates the next rotation time from now.
// rotation_hour is in the admin's local timezone (defined by timezoneOffset).
// The result is always in UTC.
func computeNextRotationAt(rotationType string, rotationDay *int32, rotationHour int32, timezoneOffset int32) *time.Time {
	if rotationType == "none" {
		return nil
	}

	now := time.Now().UTC()
	utcHour := int(rotationHour) - int(timezoneOffset)

	switch rotationType {
	case "weekly":
		day := 0
		if rotationDay != nil {
			day = int(*rotationDay)
		}
		// Find next occurrence of this weekday at the specified hour (UTC)
		daysUntil := (day - int(now.Weekday()) + 7) % 7
		candidate := time.Date(now.Year(), now.Month(), now.Day()+daysUntil, utcHour, 0, 0, 0, time.UTC)
		if !candidate.After(now) {
			candidate = candidate.AddDate(0, 0, 7)
		}
		return &candidate

	case "monthly":
		day := 1
		if rotationDay != nil {
			day = int(*rotationDay)
		}
		candidate := clampedMonthDay(now.Year(), now.Month(), day, utcHour)
		if !candidate.After(now) {
			candidate = clampedMonthDay(now.Year(), now.Month()+1, day, utcHour)
		}
		return &candidate
	}

	return nil
}

// clampedMonthDay creates a time for the given year/month/day, clamping the day
// to the last day of the month if it exceeds it.
func clampedMonthDay(year int, month time.Month, day int, hour int) time.Time {
	// Normalize month (e.g. month=13 → next year January)
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(t.Year(), t.Month(), day, hour, 0, 0, 0, time.UTC)
}

type CreatePackPoolParam struct {
	Name              string
	Description       *string
	Image             *string
	BannerType        string
	SortOrder         int32
	IsActive          bool
	OpenAt            *time.Time
	CloseAt           *time.Time
	ActiveCount       int32
	RotationType      string
	RotationDay       *int32
	RotationInterval  int32
	RotationHour      int32
	TimezoneOffset    int32
	RotationOrderMode string
	PreviewDays       int32
}

type FindAllPackPoolsParam struct {
	Search     string
	BannerType string
	ActiveOnly bool
	Page       int64
	Limit      int64
}

func (s *PackPoolService) CreateOne(ctx context.Context, param CreatePackPoolParam) (PackPoolPayload, code.I) {
	nextRotationAt := computeNextRotationAt(param.RotationType, param.RotationDay, param.RotationHour, param.TimezoneOffset)

	pool, err := s.service.repository.PackPool.CreateOne(ctx, poolRepo.CreateOneParam{
		Name:              param.Name,
		Description:       param.Description,
		ActiveCount:       param.ActiveCount,
		RotationDay:       param.RotationDay,
		Image:             param.Image,
		BannerType:        param.BannerType,
		SortOrder:         param.SortOrder,
		IsActive:          param.IsActive,
		OpenAt:            param.OpenAt,
		CloseAt:           param.CloseAt,
		RotationType:      param.RotationType,
		RotationInterval:  param.RotationInterval,
		RotationHour:      param.RotationHour,
		RotationOrderMode: param.RotationOrderMode,
		NextRotationAt:    nextRotationAt,
		PreviewDays:       param.PreviewDays,
	})
	if err != nil {
		return PackPoolPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toPackPoolPayload(pool), code.OK()
}

func (s *PackPoolService) FindOneByIdWithPacks(ctx context.Context, id string) (PackPoolPayload, code.I) {
	pool, err := s.service.repository.PackPool.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return PackPoolPayload{}, code.ModelNotFound.Err()
		}
		return PackPoolPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := toPackPoolPayload(pool)

	// Fetch all packs in this pool
	packs, err := s.service.repository.Pack.FindCurrentByPool(ctx, pool.ID.String(), 0)
	if err == nil {
		payload.Packs = toPackPoolPackItems(packs)
	}

	return payload, code.OK()
}

func (s *PackPoolService) FindAll(ctx context.Context, param FindAllPackPoolsParam) ([]PackPoolPayload, pagination.Meta, code.I) {
	offset := (param.Page - 1) * param.Limit

	total, err := s.service.repository.PackPool.Count(ctx, poolRepo.CountParam{
		Search:     param.Search,
		BannerType: param.BannerType,
		ActiveOnly: param.ActiveOnly,
	})
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	pools, err := s.service.repository.PackPool.FindAll(ctx, poolRepo.FindAllParam{
		Search:     param.Search,
		BannerType: param.BannerType,
		ActiveOnly: param.ActiveOnly,
		Limit:      param.Limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]PackPoolPayload, len(pools))
	for i, p := range pools {
		payload[i] = toPackPoolPayload(p)
	}

	return payload, pagination.NewMeta(total, param.Page, param.Limit), code.OK()
}

func (s *PackPoolService) FindActiveBanners(ctx context.Context) ([]PackPoolPayload, code.I) {
	pools, err := s.service.repository.PackPool.FindActiveBanners(ctx)
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payloads := make([]PackPoolPayload, len(pools))

	for i, pool := range pools {
		payload := toPackPoolPayload(pool)
		if string(pool.RotationType) != "none" && pool.PreviewDays > 0 && pool.NextRotationAt != nil {
			previewStart := pool.NextRotationAt.AddDate(0, 0, -int(pool.PreviewDays))
			payload.HasPreview = time.Now().UTC().After(previewStart)
		}
		packs, err := s.service.repository.Pack.FindCurrentByPool(ctx, pool.ID.String(), pool.ActiveCount)
		if err == nil {
			payload.Packs = toPackPoolPackItems(packs)
		}
		payloads[i] = payload
	}

	return payloads, code.OK()
}

func (s *PackPoolService) FindNextPacks(ctx context.Context, id string) ([]PackPoolPackItem, code.I) {
	pool, err := s.service.repository.PackPool.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, code.ModelNotFound.Err()
		}
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	if string(pool.RotationType) == "none" {
		return nil, code.OK()
	}

	nextPacks, err := s.service.repository.Pack.FindNextRotationByPool(ctx, pool.ID.String(), pool.ActiveCount)
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	items := make([]PackPoolPackItem, len(nextPacks))
	for i, p := range nextPacks {
		items[i] = PackPoolPackItem{
			ID:            p.ID,
			Code:          p.Code,
			Name:          p.Name,
			Description:   p.Description,
			Image:         p.Image,
			NameImage:     p.NameImage,
			CardsPerPull:  p.CardsPerPull,
			SortOrder:     p.SortOrder,
			Config:        unmarshalPackConfig(p.Config),
			PoolId:        p.PoolID,
			RotationOrder: p.RotationOrder,
		}
	}

	return items, code.OK()
}

func (s *PackPoolService) UpdateOneById(ctx context.Context, id string, updates map[string]any) (PackPoolPayload, code.I) {
	// Recompute next_rotation_at if rotation-related fields change
	_, hasRotationType := updates["rotation_type"]
	_, hasRotationDay := updates["rotation_day"]
	_, hasRotationInterval := updates["rotation_interval"]
	_, hasRotationHour := updates["rotation_hour"]
	_, hasTimezoneOffset := updates["timezone_offset"]

	if hasRotationType || hasRotationDay || hasRotationInterval || hasRotationHour || hasTimezoneOffset {
		// Fetch current pool to merge with updates
		current, err := s.service.repository.PackPool.FindOneById(ctx, id)
		if err != nil {
			if errors.Is(err, qrm.ErrNoRows) {
				return PackPoolPayload{}, code.ModelNotFound.Err()
			}
			return PackPoolPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}

		// Use updated values if present, otherwise use current values
		rotationType := string(current.RotationType)
		if v, ok := updates["rotation_type"].(string); ok {
			rotationType = v
		}

		rotationDay := current.RotationDay
		if v, ok := updates["rotation_day"]; ok {
			switch val := v.(type) {
			case *int32:
				rotationDay = val
			case nil:
				rotationDay = nil
			}
		}

		rotationHour := current.RotationHour
		if v, ok := updates["rotation_hour"]; ok {
			switch val := v.(type) {
			case int32:
				rotationHour = val
			case float64:
				rotationHour = int32(val)
			}
		}

		var timezoneOffset int32
		if v, ok := updates["timezone_offset"]; ok {
			switch val := v.(type) {
			case int32:
				timezoneOffset = val
			case float64:
				timezoneOffset = int32(val)
			}
		}

		nextRotationAt := computeNextRotationAt(rotationType, rotationDay, rotationHour, timezoneOffset)
		updates["next_rotation_at"] = nextRotationAt
	}

	// timezone_offset is input-only, not a DB column
	delete(updates, "timezone_offset")

	pool, err := s.service.repository.PackPool.UpdateOneById(ctx, id, updates)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return PackPoolPayload{}, code.ModelNotFound.Err()
		}
		return PackPoolPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toPackPoolPayload(pool), code.OK()
}

func (s *PackPoolService) DeleteOneById(ctx context.Context, id string) code.I {
	err := s.service.repository.PackPool.DeleteOneById(ctx, id)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}

// ReorderPacks reorders packs within a pool. ids must contain ALL pack IDs belonging to the pool.
func (s *PackPoolService) ReorderPacks(ctx context.Context, poolId string, ids []string) code.I {
	// Validate pool exists
	_, err := s.service.repository.PackPool.FindOneById(ctx, poolId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return code.ModelNotFound.Err()
		}
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	// Get all pack IDs in this pool
	existingIds, err := s.service.repository.Pack.FindIdsByPoolId(ctx, poolId)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	// Validate: must match exactly
	if len(ids) != len(existingIds) {
		return code.HttpBadRequest.Err()
	}

	existingSet := make(map[string]bool, len(existingIds))
	for _, id := range existingIds {
		existingSet[id.String()] = true
	}
	for _, id := range ids {
		if !existingSet[id] {
			return code.HttpBadRequest.Err()
		}
	}

	if err := s.service.repository.Pack.Reorder(ctx, ids); err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}

// Reorder reorders pack pools within a banner type. ids must contain ALL pool IDs of that banner type.
func (s *PackPoolService) Reorder(ctx context.Context, bannerType string, ids []string) code.I {
	existingIds, err := s.service.repository.PackPool.FindIdsByBannerType(ctx, bannerType)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	if len(ids) != len(existingIds) {
		return code.HttpBadRequest.Err()
	}

	existingSet := make(map[string]bool, len(existingIds))
	for _, id := range existingIds {
		existingSet[id.String()] = true
	}
	for _, id := range ids {
		if !existingSet[id] {
			return code.HttpBadRequest.Err()
		}
	}

	if err := s.service.repository.PackPool.Reorder(ctx, ids); err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}

// ReorderRotation reorders rotation_order of packs within a pool. ids must contain ALL pack IDs belonging to the pool.
func (s *PackPoolService) ReorderRotation(ctx context.Context, poolId string, ids []string) code.I {
	_, err := s.service.repository.PackPool.FindOneById(ctx, poolId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return code.ModelNotFound.Err()
		}
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	existingIds, err := s.service.repository.Pack.FindIdsByPoolId(ctx, poolId)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	if len(ids) != len(existingIds) {
		return code.HttpBadRequest.Err()
	}

	existingSet := make(map[string]bool, len(existingIds))
	for _, id := range existingIds {
		existingSet[id.String()] = true
	}
	for _, id := range ids {
		if !existingSet[id] {
			return code.HttpBadRequest.Err()
		}
	}

	if err := s.service.repository.Pack.ReorderRotation(ctx, ids); err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}

// AssignPacks sets the given pack IDs to this pool and unassigns any packs currently
// in the pool but not in the list. An empty list unassigns all packs.
func (s *PackPoolService) AssignPacks(ctx context.Context, poolId string, ids []string) code.I {
	// Validate pool exists
	_, err := s.service.repository.PackPool.FindOneById(ctx, poolId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return code.ModelNotFound.Err()
		}
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	if err := s.service.repository.Pack.AssignToPool(ctx, poolId, ids); err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
