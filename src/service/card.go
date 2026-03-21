package service

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zakiverse/zakiverse-api/core/code"
	cardRepo "github.com/zakiverse/zakiverse-api/src/repository/card"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CardService struct {
	service *Service
}

type CreateCardParam struct {
	MalId    int32
	AnimeId  string
	RarityId string
	Name     string
	Image    string
	Config   string
}

type CardPayload struct {
	Id       uuid.UUID `json:"id"`
	MalId    int32     `json:"mal_id"`
	AnimeId  string    `json:"anime_id"`
	RarityId string    `json:"rarity_id"`
	Name     string    `json:"name"`
	Image    string    `json:"image"`
	Config   string    `json:"config"`
}

func (s *CardService) CreateOne(ctx context.Context, param CreateCardParam) (CardPayload, code.I) {
	card, err := s.service.repository.Card.CreateOne(ctx, cardRepo.CreateOneParam{
		MalId:    param.MalId,
		AnimeId:  param.AnimeId,
		RarityId: param.RarityId,
		Name:     param.Name,
		Image:    param.Image,
		Config:   param.Config,
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) =="23505" {
			return CardPayload{}, code.CardAlreadyExists.Err()
		}
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return CardPayload{
		Id:       card.ID,
		MalId:    card.MalID,
		AnimeId:  card.AnimeID.String(),
		RarityId: card.RarityID.String(),
		Name:     card.Name,
		Image:    card.Image,
		Config:   card.Config,
	}, code.OK()
}

func (s *CardService) FindOneById(ctx context.Context, id string) (CardPayload, code.I) {
	card, err := s.service.repository.Card.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return CardPayload{}, code.ModelNotFound.Err()
		}
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return CardPayload{
		Id:       card.ID,
		MalId:    card.MalID,
		AnimeId:  card.AnimeID.String(),
		RarityId: card.RarityID.String(),
		Name:     card.Name,
		Image:    card.Image,
		Config:   card.Config,
	}, code.OK()
}

type FindAllCardsByAnimeIdParam struct {
	AnimeId string
	Page    int64
	Limit   int64
}

func (s *CardService) FindAllByAnimeId(ctx context.Context, param FindAllCardsByAnimeIdParam) ([]CardPayload, code.I) {
	offset := (param.Page - 1) * param.Limit

	cards, err := s.service.repository.Card.FindAllByAnimeId(ctx, cardRepo.FindAllByAnimeIdParam{
		AnimeId: param.AnimeId,
		Limit:   param.Limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]CardPayload, len(cards))
	for i, c := range cards {
		payload[i] = CardPayload{
			Id:       c.ID,
			MalId:    c.MalID,
			AnimeId:  c.AnimeID.String(),
			RarityId: c.RarityID.String(),
			Name:     c.Name,
			Image:    c.Image,
			Config:   c.Config,
		}
	}

	return payload, code.OK()
}

type UpdateCardParam struct {
	RarityId string
	Name     string
	Image    string
	Config   string
}

func (s *CardService) UpdateOneById(ctx context.Context, id string, param UpdateCardParam) (CardPayload, code.I) {
	card, err := s.service.repository.Card.UpdateOneById(ctx, id, cardRepo.UpdateOneByIdParam{
		RarityId: param.RarityId,
		Name:     param.Name,
		Image:    param.Image,
		Config:   param.Config,
	})
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return CardPayload{}, code.ModelNotFound.Err()
		}
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return CardPayload{
		Id:       card.ID,
		MalId:    card.MalID,
		AnimeId:  card.AnimeID.String(),
		RarityId: card.RarityID.String(),
		Name:     card.Name,
		Image:    card.Image,
		Config:   card.Config,
	}, code.OK()
}

func (s *CardService) DeleteOneById(ctx context.Context, id string) code.I {
	err := s.service.repository.Card.DeleteOneById(ctx, id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) =="23503" {
			return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
