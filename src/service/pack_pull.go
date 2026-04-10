package service

import (
	"context"
	"math/rand"

	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	historyRepo "github.com/zakiverse/zakiverse-api/src/repository/account_pull_history"
	pityRepo "github.com/zakiverse/zakiverse-api/src/repository/account_pack_pity"
	"github.com/zakiverse/zakiverse-api/src/repository/pack"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type PullResultPayload struct {
	Cards []PulledCardPayload `json:"cards"`
}

type PulledCardPayload struct {
	CardId     uuid.UUID `json:"card_id"`
	Rarity     string    `json:"rarity"`
	IsNew      bool      `json:"is_new"`
	IsPity     bool      `json:"is_pity"`
	IsFeatured bool      `json:"is_featured"`
}

func (s *PackService) Pull(ctx context.Context, accountId string, packId string, mode string) (PullResultPayload, code.I) {
	// 1. Get pack with config
	packData, err := s.service.repository.Pack.FindOneById(ctx, packId)
	if err != nil {
		return PullResultPayload{}, code.ModelNotFound.Err()
	}

	// 1a. Pool-based validation
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

		pulledCards[i] = PulledCardPayload{
			CardId:     card.CardID,
			Rarity:     rarity,
			IsPity:     isPity,
			IsFeatured: isFeatured,
		}

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

	// 7. Grant cards to account (skip duplicates)
	for i, pulled := range pulledCards {
		_, codeErr := s.service.AccountCard.AddCard(ctx, AddCardParam{
			AccountId: accountId,
			CardId:    pulled.CardId.String(),
		})
		if codeErr.OK() {
			pulledCards[i].IsNew = true
		}
		// If already owned, IsNew stays false — no error
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
		// The cards are already granted
	}

	return PullResultPayload{Cards: pulledCards}, code.OK()
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
