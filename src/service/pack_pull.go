package service

import (
	"context"
	"errors"
	"math/rand"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	accountCardRepo "github.com/zakiverse/zakiverse-api/src/repository/account_card"
	historyRepo "github.com/zakiverse/zakiverse-api/src/repository/account_pull_history"
	pityRepo "github.com/zakiverse/zakiverse-api/src/repository/account_pack_pity"
	"github.com/zakiverse/zakiverse-api/src/repository/pack"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type PullResultPayload struct {
	Cards   []PulledCardPayload `json:"cards"`
	Balance PullBalancePayload  `json:"balance"`
}

type PullBalancePayload struct {
	Coin       int32 `json:"coin"`
	CoinSpent  int32 `json:"coin_spent"`
	CoinGained int32 `json:"coin_gained"`
}

type LevelUpPayload struct {
	OldLevel int32 `json:"old_level"`
	NewLevel int32 `json:"new_level"`
}

type PulledCardPayload struct {
	CardId      uuid.UUID            `json:"card_id"`
	Rarity      string               `json:"rarity"`
	IsNew       bool                 `json:"is_new"`
	IsPity      bool                 `json:"is_pity"`
	IsFeatured  bool                 `json:"is_featured"`
	Name        string               `json:"name"`
	Image       string               `json:"image"`
	Config      CardConfig           `json:"config"`
	TagName     string               `json:"tag_name"`
	Favorite    int32                `json:"favorite"`
	Anime       PackCardAnimePayload `json:"anime"`
	LevelUp     *LevelUpPayload      `json:"level_up"`
	CoinsGained *int32               `json:"coins_gained"`
}

func (s *PackService) Pull(ctx context.Context, accountId string, packId string, mode string) (PullResultPayload, code.I) {
	// 0. Ensure account balance exists
	codeErr := s.service.AccountBalance.EnsureExists(ctx, accountId)
	if !codeErr.OK() {
		return PullResultPayload{}, codeErr
	}

	// 0a. Deduct coins
	cost := int32(s.service.config.Game.PullCostSingle)
	if mode == "multi" {
		cost = int32(s.service.config.Game.PullCostMulti)
	}
	_, err := s.service.repository.AccountBalance.DeductCoins(ctx, accountId, int(cost))
	if err != nil {
		return PullResultPayload{}, code.InsufficientCoins.Err()
	}

	// 1. Get pack with config
	packData, err := s.service.repository.Pack.FindOneById(ctx, packId)
	if err != nil {
		return PullResultPayload{}, code.ModelNotFound.Err()
	}

	// 1a. Pool-based validation
	if packData.Pack.PoolID == nil {
		return PullResultPayload{}, code.BannerNotActive.Err()
	}

	pool, err := s.service.repository.PackPool.FindOneById(ctx, packData.Pack.PoolID.String())
	if err != nil {
		return PullResultPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	if !pool.IsActive {
		return PullResultPayload{}, code.BannerNotActive.Err()
	}

	// 1b. Rotation check — if pool rotates, ensure this pack is in current rotation
	if pool.RotationType != model.RotationType_None {
		currentPacks, err := s.service.repository.Pack.FindCurrentByPool(ctx, pool.ID.String(), pool.ActiveCount)
		if err != nil {
			return PullResultPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
		}
		inRotation := false
		for _, p := range currentPacks {
			if p.ID.String() == packId {
				inRotation = true
				break
			}
		}
		if !inRotation {
			return PullResultPayload{}, code.PackNotInRotation.Err()
		}
	}

	config := unmarshalPackConfig(packData.Pack.Config)

	// 2. Get pack cards with their rarity
	packCards, err := s.service.repository.Pack.FindCardsWithRarity(ctx, packId)
	if err != nil {
		return PullResultPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}
	// 3. Group cards by rarity and filter rates to available rarities
	cardsByRarity := groupCardsByRarity(packCards)

	availableRates := make(map[string]float64)
	for rarity, rate := range config.RarityRates {
		if len(cardsByRarity[rarity]) > 0 {
			availableRates[rarity] = rate
		}
	}
	if len(availableRates) == 0 {
		return PullResultPayload{Cards: []PulledCardPayload{}}, code.OK()
	}

	// 4. Get pity counters for this account+pack
	pityCounters, err := s.service.repository.AccountPackPity.FindAllByAccountAndPack(ctx, accountId, packId)
	if err != nil {
		return PullResultPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}
	counters := toPityMap(pityCounters)

	// 5. Pull cards
	cardsPerPull := 1
	if mode == "multi" {
		cardsPerPull = int(packData.Pack.CardsPerPull)
	}
	pulledCards := make([]PulledCardPayload, cardsPerPull)

	for i := range cardsPerPull {
		rarity, isPity := rollRarity(availableRates, config.Pity, counters)
		card, isFeatured := rollCard(cardsByRarity[rarity])

		p := PulledCardPayload{
			CardId:     card.CardID,
			Rarity:     rarity,
			IsPity:     isPity,
			IsFeatured: isFeatured,
			Name:       card.Card.Name,
			Image:      card.Card.Image,
			Config:     unmarshalCardConfig(card.Card.Config),
			Favorite:   card.Card.Favorite,
			Anime: PackCardAnimePayload{
				Title:      card.Anime.Title,
				CoverImage: card.Anime.CoverImage,
			},
		}
		if card.CardTag != nil {
			p.TagName = card.CardTag.Name
		}
		pulledCards[i] = p

		// Update pity counters: increment all, reset the one we got
		for r := range availableRates {
			if r == rarity {
				counters[r] = 0
			} else {
				counters[r]++
			}
		}
	}

	// 6. Save pity counters
	upsertParams := make([]pityRepo.UpsertParam, 0, len(counters))
	for rarity, counter := range counters {
		upsertParams = append(upsertParams, pityRepo.UpsertParam{
			AccountId: accountId,
			PackId:    packId,
			Rarity:    rarity,
			Counter:   int32(counter),
		})
	}
	if err := s.service.repository.AccountPackPity.UpsertCounters(ctx, upsertParams); err != nil {
		return PullResultPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	// 7. Grant cards to account with dupe handling (level-up / overflow)
	var totalCoinsGained int32
	for i, pulled := range pulledCards {
		existing, findErr := s.service.repository.AccountCard.FindOneByAccountIdAndCardId(ctx, accountId, pulled.CardId.String())
		if findErr != nil {
			if errors.Is(findErr, qrm.ErrNoRows) {
				// New card — insert
				_, _ = s.service.repository.AccountCard.CreateOne(ctx, accountCardRepo.CreateOneParam{
					AccountId: accountId,
					CardId:    pulled.CardId.String(),
				})
				pulledCards[i].IsNew = true
			}
		} else if existing.Level < 5 {
			// Dupe, level up
			updated, levelErr := s.service.repository.AccountCard.IncrementLevel(ctx, accountId, pulled.CardId.String())
			if levelErr == nil {
				pulledCards[i].LevelUp = &LevelUpPayload{
					OldLevel: existing.Level,
					NewLevel: updated.Level,
				}
			}
		} else {
			// Dupe, already max level — overflow to coins
			overflowCoins := int32(s.service.config.Game.OverflowCoins[pulled.Rarity])
			if overflowCoins > 0 {
				s.service.repository.AccountBalance.AddCoins(ctx, accountId, int(overflowCoins))
				pulledCards[i].CoinsGained = &overflowCoins
				totalCoinsGained += overflowCoins
			}
		}
	}

	// 8. Save pull history
	historyParams := make([]historyRepo.CreateManyParam, len(pulledCards))
	for i, pulled := range pulledCards {
		historyParams[i] = historyRepo.CreateManyParam{
			AccountId:  accountId,
			PackId:     packId,
			CardId:     pulled.CardId.String(),
			Rarity:     pulled.Rarity,
			IsPity:     pulled.IsPity,
			IsFeatured: pulled.IsFeatured,
			IsNew:      pulled.IsNew,
		}
	}
	if err := s.service.repository.AccountPullHistory.CreateMany(ctx, historyParams); err != nil {
		// Log but don't fail the pull - history is non-critical
	}

	// 9. Get final balance
	finalBalance, _ := s.service.repository.AccountBalance.FindOne(ctx, accountId)

	return PullResultPayload{
		Cards: pulledCards,
		Balance: PullBalancePayload{
			Coin:       finalBalance.Coin,
			CoinSpent:  cost,
			CoinGained: totalCoinsGained,
		},
	}, code.OK()
}

func groupCardsByRarity(cards []pack.PackCardWithRarity) map[string][]pack.PackCardWithRarity {
	grouped := make(map[string][]pack.PackCardWithRarity)
	for _, c := range cards {
		rarity := string(c.Card.Rarity)
		grouped[rarity] = append(grouped[rarity], c)
	}
	return grouped
}

func toPityMap(pityCounters []model.AccountPackPity) map[string]int {
	m := make(map[string]int)
	for _, p := range pityCounters {
		m[string(p.Rarity)] = int(p.Counter)
	}
	return m
}

// rollRarity picks a rarity using weighted random, with pity override.
// Returns the rarity string and whether pity triggered.
func rollRarity(rates map[string]float64, pity map[string]int, counters map[string]int) (string, bool) {
	// Check pity first — highest threshold rarity that hit pity wins
	for rarity, threshold := range pity {
		if threshold > 0 && counters[rarity] >= threshold {
			return rarity, true
		}
	}

	// Weighted random roll
	total := 0.0
	for _, rate := range rates {
		total += rate
	}

	roll := rand.Float64() * total
	cumulative := 0.0
	var fallback string
	for rarity, rate := range rates {
		fallback = rarity
		cumulative += rate
		if roll < cumulative {
			return rarity, false
		}
	}

	return fallback, false
}

// rollCard picks a card from the pool.
// Featured cards are checked first against their featured_rate.
// If no featured card hits, falls back to weight-based random among ALL cards.
func rollCard(cards []pack.PackCardWithRarity) (pack.PackCardWithRarity, bool) {
	// 1. Roll featured cards first — each gets an independent roll against its rate
	for _, c := range cards {
		if c.IsFeatured && c.FeaturedRate != nil && *c.FeaturedRate > 0 {
			if rand.Float64() < *c.FeaturedRate {
				return c, true
			}
		}
	}

	// 2. No featured hit — weight-based random among ALL cards
	total := 0.0
	for _, c := range cards {
		total += c.Weight
	}

	roll := rand.Float64() * total
	cumulative := 0.0
	for _, c := range cards {
		cumulative += c.Weight
		if roll < cumulative {
			return c, false
		}
	}

	return cards[len(cards)-1], false
}
