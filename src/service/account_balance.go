package service

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type AccountBalanceService struct {
	service *Service
}

type BalancePayload struct {
	Coin int32 `json:"coin"`
}

func (s *AccountBalanceService) GetBalance(ctx context.Context, accountId string) (BalancePayload, code.I) {
	balance, err := s.service.repository.AccountBalance.FindOne(ctx, accountId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return BalancePayload{Coin: 0}, code.OK()
		}
		return BalancePayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return BalancePayload{Coin: balance.Coin}, code.OK()
}

func (s *AccountBalanceService) EnsureExists(ctx context.Context, accountId string) code.I {
	_, err := s.service.repository.AccountBalance.FindOne(ctx, accountId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			_, err = s.service.repository.AccountBalance.Upsert(ctx, accountId, 0)
			if err != nil {
				return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
			}
		} else {
			return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
	}
	return code.OK()
}
