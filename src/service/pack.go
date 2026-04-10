package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	packRepo "github.com/zakiverse/zakiverse-api/src/repository/pack"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type PackService struct {
	service *Service
}

type PackConfig struct {
	RarityRates map[string]float64 `json:"rarity_rates"`
	Pity        map[string]int     `json:"pity,omitempty"`
}

type PackPayload struct {
	ID            uuid.UUID         `json:"id"`
	Code          string            `json:"code"`
	Name          string            `json:"name"`
	Description   *string           `json:"description"`
	Image         string            `json:"image"`
	NameImage     *string           `json:"name_image"`
	CardsPerPull  int32             `json:"cards_per_pull"`
	SortOrder     int32             `json:"sort_order"`
	Config        PackConfig        `json:"config"`
	PoolId        uuid.UUID         `json:"pool_id"`
	RotationOrder *int32            `json:"rotation_order"`
	TotalCards    int64             `json:"total_cards"`
	Cards         []PackCardPayload `json:"cards,omitempty"`
}

type PackCardAnimePayload struct {
	Title      string  `json:"title"`
	CoverImage *string `json:"cover_image"`
}

type PackCardPayload struct {
	ID           uuid.UUID            `json:"id"`
	CardId       uuid.UUID            `json:"card_id"`
	Weight       float64              `json:"weight"`
	IsFeatured   bool                 `json:"is_featured"`
	FeaturedRate *float64             `json:"featured_rate"`
	Name         string               `json:"name"`
	Image        string               `json:"image"`
	Rarity       model.CardRarity     `json:"rarity"`
	Anime        PackCardAnimePayload `json:"anime"`
}

func marshalPackConfig(cfg PackConfig) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func unmarshalPackConfig(raw string) PackConfig {
	var cfg PackConfig
	_ = json.Unmarshal([]byte(raw), &cfg)
	return cfg
}

func toPackPayload(pack model.Pack) PackPayload {
	return PackPayload{
		ID:            pack.ID,
		Code:          pack.Code,
		Name:          pack.Name,
		Description:   pack.Description,
		Image:         pack.Image,
		NameImage:     pack.NameImage,
		CardsPerPull:  pack.CardsPerPull,
		SortOrder:     pack.SortOrder,
		Config:        unmarshalPackConfig(pack.Config),
		PoolId:        pack.PoolID,
		RotationOrder: pack.RotationOrder,
	}
}

func toPackCardPayloads(cards []packRepo.PackCardWithCard) []PackCardPayload {
	payload := make([]PackCardPayload, len(cards))
	for i, c := range cards {
		payload[i] = PackCardPayload{
			ID:           c.ID,
			CardId:       c.CardID,
			Weight:       c.Weight,
			IsFeatured:   c.IsFeatured,
			FeaturedRate: c.FeaturedRate,
			Name:         c.Card.Name,
			Image:        c.Card.Image,
			Rarity:       c.Card.Rarity,
			Anime: PackCardAnimePayload{
				Title:      c.Card.Anime.Title,
				CoverImage: c.Card.Anime.CoverImage,
			},
		}
	}
	return payload
}

type CreatePackParam struct {
	Code          string
	Name          string
	Description   *string
	Image         string
	NameImage     *string
	CardsPerPull  int32
	SortOrder     int32
	Config        PackConfig
	PoolId        string
	RotationOrder *int32
}

func (s *PackService) CreateOne(ctx context.Context, param CreatePackParam) (PackPayload, code.I) {
	configJson, err := marshalPackConfig(param.Config)
	if err != nil {
		return PackPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	pack, err := s.service.repository.Pack.CreateOne(ctx, packRepo.CreateOneParam{
		Code:          param.Code,
		Name:          param.Name,
		Description:   param.Description,
		Image:         param.Image,
		NameImage:     param.NameImage,
		CardsPerPull:  param.CardsPerPull,
		SortOrder:     param.SortOrder,
		Config:        configJson,
		PoolId:        param.PoolId,
		RotationOrder: param.RotationOrder,
	})
	if err != nil {
		return PackPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toPackPayload(pack), code.OK()
}

func (s *PackService) FindOneById(ctx context.Context, id string) (PackPayload, code.I) {
	result, err := s.service.repository.Pack.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return PackPayload{}, code.ModelNotFound.Err()
		}
		return PackPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := toPackPayload(result.Pack)
	payload.Cards = toPackCardPayloads(result.Cards)
	payload.TotalCards = int64(len(result.Cards))

	return payload, code.OK()
}

type FindAllPacksParam struct {
	Page  int64
	Limit int64
}

func (s *PackService) FindAll(ctx context.Context, param FindAllPacksParam) ([]PackPayload, code.I) {
	offset := (param.Page - 1) * param.Limit

	packs, err := s.service.repository.Pack.FindAll(ctx, packRepo.FindAllParam{
		Limit:  param.Limit,
		Offset: offset,
	})
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]PackPayload, len(packs))
	for i, p := range packs {
		payload[i] = toPackPayload(p.Pack)
		payload[i].TotalCards = p.TotalCards
	}

	return payload, code.OK()
}

func (s *PackService) UpdateOneById(ctx context.Context, id string, updates map[string]any) (PackPayload, code.I) {
	if config, ok := updates["config"]; ok && config != nil {
		raw, err := json.Marshal(config)
		if err != nil {
			return PackPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
		updates["config"] = string(raw)
	}

	pack, err := s.service.repository.Pack.UpdateOneById(ctx, id, updates)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return PackPayload{}, code.ModelNotFound.Err()
		}
		return PackPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toPackPayload(pack), code.OK()
}

func (s *PackService) DeleteOneById(ctx context.Context, id string) code.I {
	err := s.service.repository.Pack.DeleteOneById(ctx, id)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}

type AddPackCardsParam struct {
	CardId       string
	Weight       float64
	IsFeatured   bool
	FeaturedRate *float64
}

func (s *PackService) AddCards(ctx context.Context, packId string, params []AddPackCardsParam) ([]PackCardPayload, code.I) {
	repoParams := make([]packRepo.AddCardParam, len(params))
	for i, p := range params {
		weight := p.Weight
		if weight <= 0 {
			weight = 1.0
		}
		repoParams[i] = packRepo.AddCardParam{
			PackId:       packId,
			CardId:       p.CardId,
			Weight:       weight,
			IsFeatured:   p.IsFeatured,
			FeaturedRate: p.FeaturedRate,
		}
	}

	cards, err := s.service.repository.Pack.AddCards(ctx, repoParams)
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]PackCardPayload, len(cards))
	for i, c := range cards {
		payload[i] = PackCardPayload{
			ID:           c.ID,
			CardId:       c.CardID,
			Weight:       c.Weight,
			IsFeatured:   c.IsFeatured,
			FeaturedRate: c.FeaturedRate,
		}
	}

	return payload, code.OK()
}

func (s *PackService) RemoveCards(ctx context.Context, packId string, cardIds []string) code.I {
	err := s.service.repository.Pack.RemoveCards(ctx, packId, cardIds)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
