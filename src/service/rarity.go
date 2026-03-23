package service

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zakiverse/zakiverse-api/core/code"
	rarityRepo "github.com/zakiverse/zakiverse-api/src/repository/rarity"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type RarityService struct {
	service *Service
}

type CreateRarityParam struct {
	Name   string
	Config string
}

type RarityPayload struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Config string    `json:"config"`
}

func (s *RarityService) CreateOne(ctx context.Context, param CreateRarityParam) (RarityPayload, code.I) {
	rarity, err := s.service.repository.Rarity.CreateOne(ctx, rarityRepo.CreateOneParam{
		Name:   param.Name,
		Config: param.Config,
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) == "23505" {
			return RarityPayload{}, code.RarityNameConflict.Err()
		}
		return RarityPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return RarityPayload{
		Id:     rarity.ID,
		Name:   rarity.Name,
		Config: rarity.Config,
	}, code.OK()
}

func (s *RarityService) FindOneById(ctx context.Context, id string) (RarityPayload, code.I) {
	rarity, err := s.service.repository.Rarity.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return RarityPayload{}, code.ModelNotFound.Err()
		}
		return RarityPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return RarityPayload{
		Id:     rarity.ID,
		Name:   rarity.Name,
		Config: rarity.Config,
	}, code.OK()
}

func (s *RarityService) FindAll(ctx context.Context) ([]RarityPayload, code.I) {
	rarities, err := s.service.repository.Rarity.FindAll(ctx)
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]RarityPayload, len(rarities))
	for i, r := range rarities {
		payload[i] = RarityPayload{
			Id:     r.ID,
			Name:   r.Name,
			Config: r.Config,
		}
	}

	return payload, code.OK()
}

type UpdateRarityParam struct {
	Name   string
	Config string
}

func (s *RarityService) UpdateOneById(ctx context.Context, id string, param UpdateRarityParam) (RarityPayload, code.I) {
	rarity, err := s.service.repository.Rarity.UpdateOneById(ctx, id, rarityRepo.UpdateOneByIdParam{
		Name:   param.Name,
		Config: param.Config,
	})
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return RarityPayload{}, code.ModelNotFound.Err()
		}
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) == "23505" {
			return RarityPayload{}, code.RarityNameConflict.Err()
		}
		return RarityPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return RarityPayload{
		Id:     rarity.ID,
		Name:   rarity.Name,
		Config: rarity.Config,
	}, code.OK()
}

func (s *RarityService) DeleteOneById(ctx context.Context, id string) code.I {
	err := s.service.repository.Rarity.DeleteOneById(ctx, id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) == "23503" {
			return code.RarityDeleteHasCards.Err()
		}
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
