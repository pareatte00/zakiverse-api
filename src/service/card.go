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
	MalId   int32
	AnimeId string
	Rarity  string
	Name    string
	Image   string
}

type CardPayload struct {
	Id      uuid.UUID `json:"id"`
	MalId   int32     `json:"mal_id"`
	AnimeId string    `json:"anime_id"`
	Rarity  string    `json:"rarity"`
	Name    string    `json:"name"`
	Image   string    `json:"image"`
}

func (s *CardService) CreateOne(ctx context.Context, param CreateCardParam) (CardPayload, code.I) {
	card, err := s.service.repository.Card.CreateOne(ctx, cardRepo.CreateOneParam{
		MalId:   param.MalId,
		AnimeId: param.AnimeId,
		Rarity:  param.Rarity,
		Name:    param.Name,
		Image:   param.Image,
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) == "23505" {
			return CardPayload{}, code.CardAlreadyExists.Err()
		}
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return CardPayload{
		Id:      card.ID,
		MalId:   card.MalID,
		AnimeId: card.AnimeID.String(),
		Rarity:  string(card.Rarity),
		Name:    card.Name,
		Image:   card.Image,
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
		Id:      card.ID,
		MalId:   card.MalID,
		AnimeId: card.AnimeID.String(),
		Rarity:  string(card.Rarity),
		Name:    card.Name,
		Image:   card.Image,
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
			Id:      c.ID,
			MalId:   c.MalID,
			AnimeId: c.AnimeID.String(),
			Rarity:  string(c.Rarity),
			Name:    c.Name,
			Image:   c.Image,
		}
	}

	return payload, code.OK()
}

type UpdateCardParam struct {
	Rarity string
	Name   string
	Image  string
}

func (s *CardService) UpdateOneById(ctx context.Context, id string, param UpdateCardParam) (CardPayload, code.I) {
	card, err := s.service.repository.Card.UpdateOneById(ctx, id, cardRepo.UpdateOneByIdParam{
		Rarity: param.Rarity,
		Name:   param.Name,
		Image:  param.Image,
	})
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return CardPayload{}, code.ModelNotFound.Err()
		}
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return CardPayload{
		Id:      card.ID,
		MalId:   card.MalID,
		AnimeId: card.AnimeID.String(),
		Rarity:  string(card.Rarity),
		Name:    card.Name,
		Image:   card.Image,
	}, code.OK()
}

func (s *CardService) DeleteOneById(ctx context.Context, id string) code.I {
	err := s.service.repository.Card.DeleteOneById(ctx, id)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
