package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/code"
	tagRepo "github.com/zakiverse/zakiverse-api/src/repository/card_tag"
	"github.com/zakiverse/zakiverse-api/util/pagination"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CardTagService struct {
	service *Service
}

type CardTagPayload struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCardTagParam struct {
	Name string
}

type FindAllCardTagsParam struct {
	Page  int64
	Limit int64
}

var multiSpaceRe = regexp.MustCompile(`\s{2,}`)

func sanitizeTagName(name string) string {
	name = strings.TrimSpace(name)
	name = multiSpaceRe.ReplaceAllString(name, " ")
	name = strings.ToLower(name)
	return name
}

func (s *CardTagService) CreateOne(ctx context.Context, param CreateCardTagParam) (CardTagPayload, code.I) {
	name := sanitizeTagName(param.Name)

	tag, err := s.service.repository.CardTag.CreateOne(ctx, tagRepo.CreateOneParam{
		Name: name,
	})
	if err != nil {
		return CardTagPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return CardTagPayload{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}, code.OK()
}

func (s *CardTagService) FindOneById(ctx context.Context, id string) (CardTagPayload, code.I) {
	tag, err := s.service.repository.CardTag.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return CardTagPayload{}, code.ModelNotFound.Err()
		}
		return CardTagPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return CardTagPayload{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}, code.OK()
}

func (s *CardTagService) FindAll(ctx context.Context, param FindAllCardTagsParam) ([]CardTagPayload, pagination.Meta, code.I) {
	offset := (param.Page - 1) * param.Limit

	total, err := s.service.repository.CardTag.Count(ctx)
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	tags, err := s.service.repository.CardTag.FindAll(ctx, tagRepo.FindAllParam{
		Limit:  param.Limit,
		Offset: offset,
	})
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]CardTagPayload, len(tags))
	for i, t := range tags {
		payload[i] = CardTagPayload{
			ID:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		}
	}

	return payload, pagination.NewMeta(total, param.Page, param.Limit), code.OK()
}

func (s *CardTagService) UpdateOneById(ctx context.Context, id string, updates map[string]any) (CardTagPayload, code.I) {
	if name, ok := updates["name"]; ok {
		if nameStr, ok := name.(string); ok {
			updates["name"] = sanitizeTagName(nameStr)
		}
	}

	tag, err := s.service.repository.CardTag.UpdateOneById(ctx, id, updates)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return CardTagPayload{}, code.ModelNotFound.Err()
		}
		return CardTagPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return CardTagPayload{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}, code.OK()
}

func (s *CardTagService) DeleteOneById(ctx context.Context, id string) code.I {
	err := s.service.repository.CardTag.DeleteOneById(ctx, id)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
