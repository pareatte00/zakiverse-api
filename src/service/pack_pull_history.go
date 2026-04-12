package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/code"
	historyRepo "github.com/zakiverse/zakiverse-api/src/repository/account_pull_history"
	"github.com/zakiverse/zakiverse-api/util/pagination"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type PullHistoryPayload struct {
	ID         uuid.UUID `json:"id"`
	PackId     uuid.UUID `json:"pack_id"`
	CardId     uuid.UUID `json:"card_id"`
	Rarity     string    `json:"rarity"`
	IsPity     bool      `json:"is_pity"`
	IsFeatured bool      `json:"is_featured"`
	IsNew      bool      `json:"is_new"`
	PulledAt   time.Time `json:"pulled_at"`
}

type FindPullHistoryParam struct {
	AccountId string
	PackId    string
	Page      int64
	Limit     int64
}

func (s *PackService) GetPullHistory(ctx context.Context, param FindPullHistoryParam) ([]PullHistoryPayload, pagination.Meta, code.I) {
	offset := (param.Page - 1) * param.Limit

	total, err := s.service.repository.AccountPullHistory.CountByAccountAndPack(ctx, param.AccountId, param.PackId)
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	records, err := s.service.repository.AccountPullHistory.FindByAccountAndPack(ctx, historyRepo.FindByAccountAndPackParam{
		AccountId: param.AccountId,
		PackId:    param.PackId,
		Limit:     param.Limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]PullHistoryPayload, len(records))
	for i, r := range records {
		payload[i] = PullHistoryPayload{
			ID:         r.ID,
			PackId:     r.PackID,
			CardId:     r.CardID,
			Rarity:     string(r.Rarity),
			IsPity:     r.IsPity,
			IsFeatured: r.IsFeatured,
			IsNew:      r.IsNew,
			PulledAt:   r.PulledAt,
		}
	}

	return payload, pagination.NewMeta(total, param.Page, param.Limit), code.OK()
}
