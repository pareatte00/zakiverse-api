package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zakiverse/zakiverse-api/core/code"
	profileRepo "github.com/zakiverse/zakiverse-api/src/repository/profile"
	"github.com/zakiverse/zakiverse-api/util/pagination"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type ProfileService struct {
	service *Service
}

type ProfilePayload struct {
	AccountId     uuid.UUID             `json:"account_id"`
	DisplayName   string                `json:"display_name"`
	Bio           *string               `json:"bio"`
	ShowcaseCards []ShowcaseCardPayload  `json:"showcase_cards"`
	Stats         ProfileStats          `json:"stats"`
	RecentPulls   []RecentPullPayload   `json:"recent_pulls"`
}

type ShowcaseCardPayload struct {
	CardId uuid.UUID `json:"card_id"`
	Name   string    `json:"name"`
	Image  string    `json:"image"`
	Rarity string    `json:"rarity"`
	Level  int32     `json:"level"`
}

type ProfileStats struct {
	TotalCards       int64                        `json:"total_cards"`
	TotalPulls       int64                        `json:"total_pulls"`
	Completion       map[string]CompletionEntry   `json:"completion"`
	HighestLevelCard *HighestLevelCardPayload     `json:"highest_level_card"`
	LoginStreak      int32                        `json:"login_streak"`
}

type CompletionEntry struct {
	Owned   int64   `json:"owned"`
	Total   int64   `json:"total"`
	Percent float64 `json:"percent"`
}

type HighestLevelCardPayload struct {
	CardId uuid.UUID `json:"card_id"`
	Name   string    `json:"name"`
	Image  string    `json:"image"`
	Rarity string    `json:"rarity"`
	Level  int32     `json:"level"`
}

type RecentPullPayload struct {
	CardId   uuid.UUID `json:"card_id"`
	Name     string    `json:"name"`
	Image    string    `json:"image"`
	Rarity   string    `json:"rarity"`
	PulledAt string    `json:"pulled_at"`
	WasNew   bool      `json:"was_new"`
}

type ProfileSearchPayload struct {
	AccountId   uuid.UUID `json:"account_id"`
	DisplayName string    `json:"display_name"`
	TotalCards  int64     `json:"total_cards"`
	LoginStreak int32     `json:"login_streak"`
}

func (s *ProfileService) GetProfile(ctx context.Context, identifier string) (ProfilePayload, code.I) {
	// Try as UUID first, then as display name
	profile, err := s.service.repository.Profile.FindOneByAccountId(ctx, identifier)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			profile, err = s.service.repository.Profile.FindOneByDisplayName(ctx, identifier)
			if err != nil {
				if errors.Is(err, qrm.ErrNoRows) {
					return ProfilePayload{}, code.ModelNotFound.Err()
				}
				return ProfilePayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
			}
		} else {
			return ProfilePayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
	}

	accountId := profile.AccountID.String()

	// Build stats
	stats, codeErr := s.buildStats(ctx, accountId)
	if !codeErr.OK() {
		return ProfilePayload{}, codeErr
	}

	// Build showcase cards
	showcaseCards := s.buildShowcaseCards(ctx, accountId, profile.ShowcaseCards)

	// Build recent pulls
	recentPulls := s.buildRecentPulls(ctx, accountId)

	return ProfilePayload{
		AccountId:     profile.AccountID,
		DisplayName:   profile.DisplayName,
		Bio:           profile.Bio,
		ShowcaseCards: showcaseCards,
		Stats:         stats,
		RecentPulls:   recentPulls,
	}, code.OK()
}

type UpdateProfileParam struct {
	AccountId     string
	DisplayName   string
	Bio           *string
	ShowcaseCards []string
}

func (s *ProfileService) UpdateProfile(ctx context.Context, param UpdateProfileParam) (ProfilePayload, code.I) {
	// Validate showcase cards count
	if len(param.ShowcaseCards) > 6 {
		return ProfilePayload{}, code.ShowcaseTooMany.Err()
	}

	// Validate showcase cards are owned
	if len(param.ShowcaseCards) > 0 {
		cardIds := make([]uuid.UUID, len(param.ShowcaseCards))
		for i, id := range param.ShowcaseCards {
			parsed, err := uuid.Parse(id)
			if err != nil {
				return ProfilePayload{}, code.HttpBadRequest.Err()
			}
			cardIds[i] = parsed
		}
		ownedMap, err := s.service.repository.AccountCard.FindOwnedCardIds(ctx, param.AccountId, cardIds)
		if err != nil {
			return ProfilePayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
		for _, id := range cardIds {
			if !ownedMap[id] {
				return ProfilePayload{}, code.ShowcaseCardNotOwned.Err()
			}
		}
	}

	showcaseJson, err := json.Marshal(param.ShowcaseCards)
	if err != nil {
		return ProfilePayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	_, err = s.service.repository.Profile.Upsert(ctx, profileRepo.UpsertParam{
		AccountId:     param.AccountId,
		DisplayName:   param.DisplayName,
		Bio:           param.Bio,
		ShowcaseCards: string(showcaseJson),
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && string(pgErr.Code) == "23505" {
			return ProfilePayload{}, code.DisplayNameTaken.Err()
		}
		return ProfilePayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return s.GetProfile(ctx, param.AccountId)
}

func (s *ProfileService) SearchProfiles(ctx context.Context, query string, page int64, limit int64) ([]ProfileSearchPayload, pagination.Meta, code.I) {
	offset := (page - 1) * limit

	total, err := s.service.repository.Profile.CountSearch(ctx, query)
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	profiles, err := s.service.repository.Profile.Search(ctx, profileRepo.SearchParam{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]ProfileSearchPayload, len(profiles))
	for i, p := range profiles {
		totalCards, _ := s.service.repository.Profile.CountCards(ctx, p.AccountID.String())

		var loginStreak int32
		records, recErr := s.service.repository.CheckInRecord.FindByAccountId(ctx, p.AccountID.String())
		if recErr == nil {
			for _, r := range records {
				if r.Streak > loginStreak {
					loginStreak = r.Streak
				}
			}
		}

		payload[i] = ProfileSearchPayload{
			AccountId:   p.AccountID,
			DisplayName: p.DisplayName,
			TotalCards:  totalCards,
			LoginStreak: loginStreak,
		}
	}

	return payload, pagination.NewMeta(total, page, limit), code.OK()
}

func (s *ProfileService) buildStats(ctx context.Context, accountId string) (ProfileStats, code.I) {
	totalCards, err := s.service.repository.Profile.CountCards(ctx, accountId)
	if err != nil {
		return ProfileStats{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	totalPulls, err := s.service.repository.Profile.CountPulls(ctx, accountId)
	if err != nil {
		return ProfileStats{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	completionStats, err := s.service.repository.Profile.GetCompletionStats(ctx, accountId)
	if err != nil {
		return ProfileStats{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	completion := make(map[string]CompletionEntry)
	for _, cs := range completionStats {
		percent := 0.0
		if cs.Total > 0 {
			percent = float64(cs.Owned) / float64(cs.Total) * 100
			percent = float64(int(percent*10)) / 10
		}
		completion[cs.Rarity] = CompletionEntry{
			Owned:   cs.Owned,
			Total:   cs.Total,
			Percent: percent,
		}
	}

	var highestCard *HighestLevelCardPayload
	hlc, err := s.service.repository.Profile.GetHighestLevelCard(ctx, accountId)
	if err == nil && hlc != nil {
		highestCard = &HighestLevelCardPayload{
			CardId: hlc.CardId,
			Name:   hlc.Name,
			Image:  hlc.Image,
			Rarity: hlc.Rarity,
			Level:  hlc.Level,
		}
	}

	var loginStreak int32
	records, err := s.service.repository.CheckInRecord.FindByAccountId(ctx, accountId)
	if err == nil {
		for _, r := range records {
			if r.Streak > loginStreak {
				loginStreak = r.Streak
			}
		}
	}

	return ProfileStats{
		TotalCards:       totalCards,
		TotalPulls:       totalPulls,
		Completion:       completion,
		HighestLevelCard: highestCard,
		LoginStreak:      loginStreak,
	}, code.OK()
}

func (s *ProfileService) buildShowcaseCards(ctx context.Context, accountId string, showcaseJson string) []ShowcaseCardPayload {
	var cardIdStrings []string
	_ = json.Unmarshal([]byte(showcaseJson), &cardIdStrings)

	if len(cardIdStrings) == 0 {
		return []ShowcaseCardPayload{}
	}

	cardIds := make([]uuid.UUID, 0, len(cardIdStrings))
	for _, id := range cardIdStrings {
		parsed, err := uuid.Parse(id)
		if err == nil {
			cardIds = append(cardIds, parsed)
		}
	}

	if len(cardIds) == 0 {
		return []ShowcaseCardPayload{}
	}

	// Get card data + level
	ownedMap, err := s.service.repository.AccountCard.FindOwnedWithLevel(ctx, accountId, cardIds)
	if err != nil {
		return []ShowcaseCardPayload{}
	}

	result := make([]ShowcaseCardPayload, 0, len(cardIds))
	for _, id := range cardIds {
		info, owned := ownedMap[id]
		if !owned {
			continue
		}

		// Get card details
		card, err := s.service.repository.Card.FindOneById(ctx, id.String())
		if err != nil {
			continue
		}

		result = append(result, ShowcaseCardPayload{
			CardId: id,
			Name:   card.Card.Name,
			Image:  card.Card.Image,
			Rarity: string(card.Card.Rarity),
			Level:  info.Level,
		})
	}

	return result
}

func (s *ProfileService) buildRecentPulls(ctx context.Context, accountId string) []RecentPullPayload {
	pulls, err := s.service.repository.Profile.GetRecentPulls(ctx, accountId, 10)
	if err != nil {
		return []RecentPullPayload{}
	}

	result := make([]RecentPullPayload, len(pulls))
	for i, p := range pulls {
		result[i] = RecentPullPayload{
			CardId:   p.CardId,
			Name:     p.Name,
			Image:    p.Image,
			Rarity:   p.Rarity,
			PulledAt: fmt.Sprintf("%s", p.PulledAt),
			WasNew:   p.IsNew,
		}
	}

	return result
}
