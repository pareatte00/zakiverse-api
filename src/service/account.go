package service

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/code"
	discordOutbound "github.com/zakiverse/zakiverse-api/src/outbound/discord"
	accountRepo "github.com/zakiverse/zakiverse-api/src/repository/account"
	jwtutil "github.com/zakiverse/zakiverse-api/util/jwt"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type AccountService struct {
	service *Service
}

type AuthDiscordParam struct {
	Code string
}

type AuthDiscordPayload struct {
	AccessToken string `json:"access_token"`
}

func (s *AccountService) AuthDiscord(ctx context.Context, param AuthDiscordParam) (AuthDiscordPayload, code.I) {
	tokenResp, err := s.service.outbound.Discord.ExchangeCode(ctx, discordOutbound.ExchangeCodeParam{
		ClientId:     s.service.credential.DiscordClientId,
		ClientSecret: s.service.credential.DiscordClientSecret,
		RedirectURI:  s.service.config.Auth.DiscordRedirectUri,
		Code:         param.Code,
	})
	if err != nil {
		return AuthDiscordPayload{}, code.AccountDiscordAuthFailed.Err().WithError(trace.Wrap(err))
	}

	discordUser, err := s.service.outbound.Discord.GetUser(ctx, tokenResp.AccessToken)
	if err != nil {
		return AuthDiscordPayload{}, code.AccountDiscordAuthFailed.Err().WithError(trace.Wrap(err))
	}

	var avatar *string
	if discordUser.Avatar != "" {
		avatar = &discordUser.Avatar
	}

	account, err := s.service.repository.Account.FindOneByDiscordId(ctx, discordUser.ID)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return AuthDiscordPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}

		account, err = s.service.repository.Account.CreateOne(ctx, accountRepo.CreateOneParam{
			DiscordId: discordUser.ID,
			Username:  discordUser.Username,
			Email:     discordUser.Email,
			Avatar:    avatar,
		})
		if err != nil {
			return AuthDiscordPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
	} else {
		account, err = s.service.repository.Account.UpdateOneByDiscordId(ctx, discordUser.ID, map[string]any{
			"username": discordUser.Username,
			"email":    discordUser.Email,
			"avatar":   avatar,
		})
		if err != nil {
			return AuthDiscordPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
	}

	token, err := jwtutil.Generate(jwtutil.GenerateParam{
		AccountId:     account.ID,
		Username:      account.Username,
		Role:          string(account.Role),
		Secret:        s.service.credential.JwtSecret,
		ExpiryMinutes: s.service.config.Auth.AccessTokenExpiryMin,
	})
	if err != nil {
		return AuthDiscordPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return AuthDiscordPayload{
		AccessToken: token,
	}, code.OK()
}

type FindOneByIdParam struct {
	AccountId uuid.UUID
}

type FindOneByIdPayload struct {
	ID        uuid.UUID `json:"id"`
	DiscordId string    `json:"discord_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    *string   `json:"avatar"`
	Role      string    `json:"role"`
}

func (s *AccountService) FindOneById(ctx context.Context, param FindOneByIdParam) (FindOneByIdPayload, code.I) {
	account, err := s.service.repository.Account.FindOneById(ctx, param.AccountId.String())
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return FindOneByIdPayload{}, code.ModelNotFound.Err()
		}
		return FindOneByIdPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return FindOneByIdPayload{
		ID:        account.ID,
		DiscordId: account.DiscordID,
		Username:  account.Username,
		Email:     account.Email,
		Avatar:    account.Avatar,
		Role:      string(account.Role),
	}, code.OK()
}
