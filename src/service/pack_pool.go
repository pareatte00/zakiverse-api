package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/code"
	poolRepo "github.com/zakiverse/zakiverse-api/src/repository/pack_pool"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type PackPoolService struct {
	service *Service
}

type PackPoolPayload struct {
	Id            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	ActiveCount   int32      `json:"active_count"`
	RotationDay   int32      `json:"rotation_day"`
	LastRotatedAt *time.Time `json:"last_rotated_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreatePackPoolParam struct {
	Name        string
	Description *string
	ActiveCount int32
	RotationDay int32
}

type FindAllPackPoolsParam struct {
	Page  int64
	Limit int64
}

func (s *PackPoolService) CreateOne(ctx context.Context, param CreatePackPoolParam) (PackPoolPayload, code.I) {
	pool, err := s.service.repository.PackPool.CreateOne(ctx, poolRepo.CreateOneParam{
		Name:        param.Name,
		Description: param.Description,
		ActiveCount: param.ActiveCount,
		RotationDay: param.RotationDay,
	})
	if err != nil {
		return PackPoolPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return PackPoolPayload{
		Id:            pool.ID,
		Name:          pool.Name,
		Description:   pool.Description,
		ActiveCount:   pool.ActiveCount,
		RotationDay:   pool.RotationDay,
		LastRotatedAt: pool.LastRotatedAt,
		CreatedAt:     pool.CreatedAt,
		UpdatedAt:     pool.UpdatedAt,
	}, code.OK()
}

func (s *PackPoolService) FindOneById(ctx context.Context, id string) (PackPoolPayload, code.I) {
	pool, err := s.service.repository.PackPool.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return PackPoolPayload{}, code.ModelNotFound.Err()
		}
		return PackPoolPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return PackPoolPayload{
		Id:            pool.ID,
		Name:          pool.Name,
		Description:   pool.Description,
		ActiveCount:   pool.ActiveCount,
		RotationDay:   pool.RotationDay,
		LastRotatedAt: pool.LastRotatedAt,
		CreatedAt:     pool.CreatedAt,
		UpdatedAt:     pool.UpdatedAt,
	}, code.OK()
}

func (s *PackPoolService) FindAll(ctx context.Context, param FindAllPackPoolsParam) ([]PackPoolPayload, code.I) {
	offset := (param.Page - 1) * param.Limit

	pools, err := s.service.repository.PackPool.FindAll(ctx, poolRepo.FindAllParam{
		Limit:  param.Limit,
		Offset: offset,
	})
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]PackPoolPayload, len(pools))
	for i, p := range pools {
		payload[i] = PackPoolPayload{
			Id:            p.ID,
			Name:          p.Name,
			Description:   p.Description,
			ActiveCount:   p.ActiveCount,
			RotationDay:   p.RotationDay,
			LastRotatedAt: p.LastRotatedAt,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		}
	}

	return payload, code.OK()
}

func (s *PackPoolService) UpdateOneById(ctx context.Context, id string, updates map[string]any) (PackPoolPayload, code.I) {
	pool, err := s.service.repository.PackPool.UpdateOneById(ctx, id, updates)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return PackPoolPayload{}, code.ModelNotFound.Err()
		}
		return PackPoolPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return PackPoolPayload{
		Id:            pool.ID,
		Name:          pool.Name,
		Description:   pool.Description,
		ActiveCount:   pool.ActiveCount,
		RotationDay:   pool.RotationDay,
		LastRotatedAt: pool.LastRotatedAt,
		CreatedAt:     pool.CreatedAt,
		UpdatedAt:     pool.UpdatedAt,
	}, code.OK()
}

func (s *PackPoolService) DeleteOneById(ctx context.Context, id string) code.I {
	err := s.service.repository.PackPool.DeleteOneById(ctx, id)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
