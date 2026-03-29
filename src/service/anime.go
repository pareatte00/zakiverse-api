package service

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zakiverse/zakiverse-api/core/code"
	animeRepo "github.com/zakiverse/zakiverse-api/src/repository/anime"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type AnimeService struct {
	service *Service
}

type CreateAnimeParam struct {
	MalId      int32
	Title      string
	Synopsis   *string
	CoverImage *string
}

type AnimePayload struct {
	Id         uuid.UUID `json:"id"`
	MalId      int32     `json:"mal_id"`
	Title      string    `json:"title"`
	Synopsis   *string   `json:"synopsis"`
	CoverImage *string   `json:"cover_image"`
}

func (s *AnimeService) CreateOne(ctx context.Context, param CreateAnimeParam) (AnimePayload, code.I) {
	anime, err := s.service.repository.Anime.CreateOne(ctx, animeRepo.CreateOneParam{
		MalId:      param.MalId,
		Title:      param.Title,
		Synopsis:   param.Synopsis,
		CoverImage: param.CoverImage,
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) == "23505" {
			return AnimePayload{}, code.AnimeAlreadyExists.Err()
		}
		return AnimePayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return AnimePayload{
		Id:         anime.ID,
		MalId:      anime.MalID,
		Title:      anime.Title,
		Synopsis:   anime.Synopsis,
		CoverImage: anime.CoverImage,
	}, code.OK()
}

func (s *AnimeService) FindOneById(ctx context.Context, id string) (AnimePayload, code.I) {
	anime, err := s.service.repository.Anime.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return AnimePayload{}, code.ModelNotFound.Err()
		}
		return AnimePayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return AnimePayload{
		Id:         anime.ID,
		MalId:      anime.MalID,
		Title:      anime.Title,
		Synopsis:   anime.Synopsis,
		CoverImage: anime.CoverImage,
	}, code.OK()
}

type FindAllAnimeParam struct {
	Page  int64
	Limit int64
}

func (s *AnimeService) FindAll(ctx context.Context, param FindAllAnimeParam) ([]AnimePayload, code.I) {
	offset := (param.Page - 1) * param.Limit

	animes, err := s.service.repository.Anime.FindAll(ctx, animeRepo.FindAllParam{
		Limit:  param.Limit,
		Offset: offset,
	})
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]AnimePayload, len(animes))
	for i, a := range animes {
		payload[i] = AnimePayload{
			Id:         a.ID,
			MalId:      a.MalID,
			Title:      a.Title,
			Synopsis:   a.Synopsis,
			CoverImage: a.CoverImage,
		}
	}

	return payload, code.OK()
}

type UpdateAnimeParam struct {
	Title      string
	Synopsis   *string
	CoverImage *string
}

func (s *AnimeService) UpdateOneById(ctx context.Context, id string, param UpdateAnimeParam) (AnimePayload, code.I) {
	updates := map[string]any{
		"title":       param.Title,
		"synopsis":    param.Synopsis,
		"cover_image": param.CoverImage,
	}

	anime, err := s.service.repository.Anime.UpdateOneById(ctx, id, updates)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return AnimePayload{}, code.ModelNotFound.Err()
		}
		return AnimePayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return AnimePayload{
		Id:         anime.ID,
		MalId:      anime.MalID,
		Title:      anime.Title,
		Synopsis:   anime.Synopsis,
		CoverImage: anime.CoverImage,
	}, code.OK()
}

func (s *AnimeService) DeleteOneById(ctx context.Context, id string) code.I {
	err := s.service.repository.Anime.DeleteOneById(ctx, id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) == "23503" {
			return code.AnimeDeleteHasCards.Err()
		}
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
