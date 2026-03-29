package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	animeRepo "github.com/zakiverse/zakiverse-api/src/repository/anime"
	cardRepo "github.com/zakiverse/zakiverse-api/src/repository/card"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CardService struct {
	service *Service
}

type CardConfig struct {
	BackgroundImage *string `json:"background_image,omitempty"`
}

type CreateCardParam struct {
	MalId           int32
	Rarity          string
	Name            string
	Image           string
	Config          CardConfig
	AnimeMalId      int32
	AnimeTitle      string
	AnimeSynopsis   *string
	AnimeCoverImage *string
}

type CardPayload struct {
	Id     uuid.UUID    `json:"id"`
	MalId  int32        `json:"mal_id"`
	Rarity string       `json:"rarity"`
	Name   string       `json:"name"`
	Image  string       `json:"image"`
	Config CardConfig   `json:"config"`
	Anime  AnimePayload `json:"anime"`
}

func marshalCardConfig(cfg CardConfig) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func unmarshalCardConfig(raw string) CardConfig {
	var cfg CardConfig
	_ = json.Unmarshal([]byte(raw), &cfg)
	return cfg
}

func toCardPayload(card model.Card, anime model.Anime) CardPayload {
	return CardPayload{
		Id:     card.ID,
		MalId:  card.MalID,
		Rarity: string(card.Rarity),
		Name:   card.Name,
		Image:  card.Image,
		Config: unmarshalCardConfig(card.Config),
		Anime: AnimePayload{
			Id:         anime.ID,
			MalId:      anime.MalID,
			Title:      anime.Title,
			Synopsis:   anime.Synopsis,
			CoverImage: anime.CoverImage,
		},
	}
}

func (s *CardService) CreateOne(ctx context.Context, param CreateCardParam) (CardPayload, code.I) {
	anime, err := s.service.repository.Anime.FindOneByMalId(ctx, param.AnimeMalId)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}

		anime, err = s.service.repository.Anime.CreateOne(ctx, animeRepo.CreateOneParam{
			MalId:      param.AnimeMalId,
			Title:      param.AnimeTitle,
			Synopsis:   param.AnimeSynopsis,
			CoverImage: param.AnimeCoverImage,
		})
		if err != nil {
			return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
	}

	configJson, err := marshalCardConfig(param.Config)
	if err != nil {
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	card, err := s.service.repository.Card.CreateOne(ctx, cardRepo.CreateOneParam{
		MalId:   param.MalId,
		AnimeId: anime.ID.String(),
		Rarity:  param.Rarity,
		Name:    param.Name,
		Image:   param.Image,
		Config:  configJson,
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) == "23505" {
			return CardPayload{}, code.CardAlreadyExists.Err()
		}
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toCardPayload(card, anime), code.OK()
}

func (s *CardService) FindOneById(ctx context.Context, id string) (CardPayload, code.I) {
	result, err := s.service.repository.Card.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return CardPayload{}, code.ModelNotFound.Err()
		}
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toCardPayload(result.Card, result.Anime), code.OK()
}

type FindAllCardsParam struct {
	Search string
	Rarity string
	Sort   string
	Order  string
	Page   int64
	Limit  int64
}

func (s *CardService) FindAll(ctx context.Context, param FindAllCardsParam) ([]CardPayload, code.I) {
	offset := (param.Page - 1) * param.Limit

	results, err := s.service.repository.Card.FindAll(ctx, cardRepo.FindAllParam{
		Search: param.Search,
		Rarity: param.Rarity,
		Sort:   param.Sort,
		Order:  param.Order,
		Limit:  param.Limit,
		Offset: offset,
	})
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]CardPayload, len(results))
	for i, r := range results {
		payload[i] = toCardPayload(r.Card, r.Anime)
	}

	return payload, code.OK()
}

type FindAllCardsByAnimeIdParam struct {
	AnimeId string
	Page    int64
	Limit   int64
}

func (s *CardService) FindAllByAnimeId(ctx context.Context, param FindAllCardsByAnimeIdParam) ([]CardPayload, code.I) {
	offset := (param.Page - 1) * param.Limit

	results, err := s.service.repository.Card.FindAllByAnimeId(ctx, cardRepo.FindAllByAnimeIdParam{
		AnimeId: param.AnimeId,
		Limit:   param.Limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]CardPayload, len(results))
	for i, r := range results {
		payload[i] = toCardPayload(r.Card, r.Anime)
	}

	return payload, code.OK()
}

func (s *CardService) UpdateOneById(ctx context.Context, id string, updates map[string]any) (CardPayload, code.I) {
	if config, ok := updates["config"]; ok && config != nil {
		raw, err := json.Marshal(config)
		if err != nil {
			return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
		updates["config"] = string(raw)
	}

	card, err := s.service.repository.Card.UpdateOneById(ctx, id, updates)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return CardPayload{}, code.ModelNotFound.Err()
		}
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	anime, err := s.service.repository.Anime.FindOneById(ctx, card.AnimeID.String())
	if err != nil {
		return CardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toCardPayload(card, anime), code.OK()
}

func (s *CardService) DeleteOneById(ctx context.Context, id string) code.I {
	err := s.service.repository.Card.DeleteOneById(ctx, id)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
