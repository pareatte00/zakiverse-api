package service

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zakiverse/zakiverse-api/core/code"
	accountCardRepo "github.com/zakiverse/zakiverse-api/src/repository/account_card"
	"github.com/zakiverse/zakiverse-api/util/pagination"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type AccountCardService struct {
	service *Service
}

type AccountCardPayload struct {
	ID         uuid.UUID `json:"id"`
	AccountId  string    `json:"account_id"`
	CardId     string    `json:"card_id"`
	ObtainedAt string    `json:"obtained_at"`
}

type FindMyCardsParam struct {
	AccountId string
	Page      int64
	Limit     int64
}

func (s *AccountCardService) FindMyCards(ctx context.Context, param FindMyCardsParam) ([]AccountCardPayload, pagination.Meta, code.I) {
	offset := (param.Page - 1) * param.Limit

	total, err := s.service.repository.AccountCard.CountByAccountId(ctx, param.AccountId)
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	accountCards, err := s.service.repository.AccountCard.FindAllByAccountId(ctx, accountCardRepo.FindAllByAccountIdParam{
		AccountId: param.AccountId,
		Limit:     param.Limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]AccountCardPayload, len(accountCards))
	for i, ac := range accountCards {
		payload[i] = AccountCardPayload{
			ID:         ac.ID,
			AccountId:  ac.AccountID.String(),
			CardId:     ac.CardID.String(),
			ObtainedAt: ac.ObtainedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return payload, pagination.NewMeta(total, param.Page, param.Limit), code.OK()
}

type AddCardParam struct {
	AccountId string
	CardId    string
}

func (s *AccountCardService) AddCard(ctx context.Context, param AddCardParam) (AccountCardPayload, code.I) {
	accountCard, err := s.service.repository.AccountCard.CreateOne(ctx, accountCardRepo.CreateOneParam{
		AccountId: param.AccountId,
		CardId:    param.CardId,
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) == "23505" {
			return AccountCardPayload{}, code.AccountCardAlreadyOwned.Err()
		}
		return AccountCardPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return AccountCardPayload{
		ID:         accountCard.ID,
		AccountId:  accountCard.AccountID.String(),
		CardId:     accountCard.CardID.String(),
		ObtainedAt: accountCard.ObtainedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, code.OK()
}

type RemoveCardParam struct {
	AccountId string
	CardId    string
}

func (s *AccountCardService) RemoveCard(ctx context.Context, param RemoveCardParam) code.I {
	_, err := s.service.repository.AccountCard.FindOneByAccountIdAndCardId(ctx, param.AccountId, param.CardId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return code.ModelNotFound.Err()
		}
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	err = s.service.repository.AccountCard.DeleteOne(ctx, param.AccountId, param.CardId)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}
