package service

import (
	"context"

	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type PityInfoPayload struct {
	Rarity    string `json:"rarity"`
	Counter   int    `json:"counter"`
	Threshold int    `json:"threshold"`
}

func (s *PackService) GetPityInfo(ctx context.Context, accountId string, packId string) ([]PityInfoPayload, code.I) {
	// 1. Get pack config for pity thresholds
	packData, err := s.service.repository.Pack.FindOneById(ctx, packId)
	if err != nil {
		return nil, code.ModelNotFound.Err()
	}

	config := unmarshalPackConfig(packData.Pack.Config)

	// 2. Get current pity counters
	pityCounters, err := s.service.repository.AccountPackPity.FindAllByAccountAndPack(ctx, accountId, packId)
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}
	counters := toPityMap(pityCounters)

	// 3. Build payload for each rarity that has a pity threshold
	payload := make([]PityInfoPayload, 0, len(config.Pity))
	for rarity, threshold := range config.Pity {
		if threshold <= 0 {
			continue
		}
		payload = append(payload, PityInfoPayload{
			Rarity:    rarity,
			Counter:   counters[rarity],
			Threshold: threshold,
		})
	}

	return payload, code.OK()
}
