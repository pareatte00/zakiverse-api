package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	checkInPlanRepo "github.com/zakiverse/zakiverse-api/src/repository/check_in_plan"
	checkInRecordRepo "github.com/zakiverse/zakiverse-api/src/repository/check_in_record"
	"github.com/zakiverse/zakiverse-api/util/pagination"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type CheckInService struct {
	service *Service
}

type CheckInPlanPayload struct {
	ID          uuid.UUID        `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Type        string           `json:"type"`
	Interval    int32            `json:"interval"`
	MaxClaims   int32            `json:"max_claims"`
	Rewards     json.RawMessage  `json:"rewards"`
	ResetPolicy string           `json:"reset_policy"`
	IsActive    bool             `json:"is_active"`
	StartsAt    *time.Time       `json:"starts_at"`
	EndsAt      *time.Time       `json:"ends_at"`
	SortOrder   int32            `json:"sort_order"`
	Status      *CheckInStatus   `json:"status,omitempty"`
}

type CheckInStatus struct {
	CanClaim    bool       `json:"can_claim"`
	NextClaimAt *time.Time `json:"next_claim_at"`
	Streak      int32      `json:"streak"`
	ClaimCount  int32      `json:"claim_count"`
}

type CheckInReward struct {
	Claim  string `json:"claim,omitempty"`
	Day    int    `json:"day,omitempty"`
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

type ClaimResultPayload struct {
	Reward  CheckInReward  `json:"reward"`
	Status  CheckInStatus  `json:"status"`
	Balance BalancePayload `json:"balance"`
}

func toPlanPayload(plan model.CheckInPlan) CheckInPlanPayload {
	return CheckInPlanPayload{
		ID:          plan.ID,
		Code:        plan.Code,
		Name:        plan.Name,
		Description: plan.Description,
		Type:        string(plan.Type),
		Interval:    plan.Interval,
		MaxClaims:   plan.MaxClaims,
		Rewards:     json.RawMessage(plan.Rewards),
		ResetPolicy: string(plan.ResetPolicy),
		IsActive:    plan.IsActive,
		StartsAt:    plan.StartsAt,
		EndsAt:      plan.EndsAt,
		SortOrder:   plan.SortOrder,
	}
}

// --- Admin CRUD ---

type CreateCheckInPlanParam struct {
	Code        string
	Name        string
	Description *string
	Type        string
	Interval    int32
	MaxClaims   int32
	Rewards     string
	ResetPolicy string
	IsActive    bool
	StartsAt    *time.Time
	EndsAt      *time.Time
	SortOrder   int32
}

func (s *CheckInService) CreatePlan(ctx context.Context, param CreateCheckInPlanParam) (CheckInPlanPayload, code.I) {
	plan, err := s.service.repository.CheckInPlan.CreateOne(ctx, checkInPlanRepo.CreateOneParam{
		Code:        param.Code,
		Name:        param.Name,
		Description: param.Description,
		Type:        param.Type,
		Interval:    param.Interval,
		MaxClaims:   param.MaxClaims,
		Rewards:     param.Rewards,
		ResetPolicy: param.ResetPolicy,
		IsActive:    param.IsActive,
		StartsAt:    param.StartsAt,
		EndsAt:      param.EndsAt,
		SortOrder:   param.SortOrder,
	})
	if err != nil {
		return CheckInPlanPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toPlanPayload(plan), code.OK()
}

type FindAllPlansParam struct {
	Page  int64
	Limit int64
}

func (s *CheckInService) FindAllPlans(ctx context.Context, param FindAllPlansParam) ([]CheckInPlanPayload, pagination.Meta, code.I) {
	offset := (param.Page - 1) * param.Limit

	total, err := s.service.repository.CheckInPlan.Count(ctx)
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	plans, err := s.service.repository.CheckInPlan.FindAll(ctx, checkInPlanRepo.FindAllParam{
		Limit:  param.Limit,
		Offset: offset,
	})
	if err != nil {
		return nil, pagination.Meta{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	payload := make([]CheckInPlanPayload, len(plans))
	for i, p := range plans {
		payload[i] = toPlanPayload(p)
	}

	return payload, pagination.NewMeta(total, param.Page, param.Limit), code.OK()
}

func (s *CheckInService) FindPlanById(ctx context.Context, id string) (CheckInPlanPayload, code.I) {
	plan, err := s.service.repository.CheckInPlan.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return CheckInPlanPayload{}, code.ModelNotFound.Err()
		}
		return CheckInPlanPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toPlanPayload(plan), code.OK()
}

func (s *CheckInService) UpdatePlanById(ctx context.Context, id string, updates map[string]any) (CheckInPlanPayload, code.I) {
	plan, err := s.service.repository.CheckInPlan.UpdateOneById(ctx, id, updates)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return CheckInPlanPayload{}, code.ModelNotFound.Err()
		}
		return CheckInPlanPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return toPlanPayload(plan), code.OK()
}

func (s *CheckInService) DeletePlanById(ctx context.Context, id string) code.I {
	err := s.service.repository.CheckInPlan.DeleteOneById(ctx, id)
	if err != nil {
		return code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	return code.OK()
}

// --- Player endpoints ---

func (s *CheckInService) GetActivePlans(ctx context.Context, accountId string) ([]CheckInPlanPayload, code.I) {
	plans, err := s.service.repository.CheckInPlan.FindActive(ctx)
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	records, err := s.service.repository.CheckInRecord.FindByAccountId(ctx, accountId)
	if err != nil {
		return nil, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	recordMap := make(map[uuid.UUID]model.CheckInRecord)
	for _, r := range records {
		recordMap[r.PlanID] = r
	}

	payload := make([]CheckInPlanPayload, len(plans))
	for i, p := range plans {
		payload[i] = toPlanPayload(p)
		status := computeStatus(p, recordMap[p.ID])
		payload[i].Status = &status
	}

	return payload, code.OK()
}

func (s *CheckInService) Claim(ctx context.Context, accountId string, planId string) (ClaimResultPayload, code.I) {
	// 1. Get plan
	plan, err := s.service.repository.CheckInPlan.FindOneById(ctx, planId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return ClaimResultPayload{}, code.ModelNotFound.Err()
		}
		return ClaimResultPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	// 2. Validate plan active + time window
	if !plan.IsActive {
		return ClaimResultPayload{}, code.CheckInPlanInactive.Err()
	}
	now := time.Now()
	if plan.StartsAt != nil && now.Before(*plan.StartsAt) {
		return ClaimResultPayload{}, code.CheckInPlanInactive.Err()
	}
	if plan.EndsAt != nil && now.After(*plan.EndsAt) {
		return ClaimResultPayload{}, code.CheckInPlanInactive.Err()
	}

	// 3. Get or create record
	record, err := s.service.repository.CheckInRecord.FindOrCreate(ctx, accountId, planId)
	if err != nil {
		return ClaimResultPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	// 4. Check interval
	if record.LastClaimed != nil {
		elapsed := now.Sub(*record.LastClaimed)
		if elapsed.Seconds() < float64(plan.Interval) {
			return ClaimResultPayload{}, code.CheckInTooEarly.Err()
		}
	}

	// 5. Check max claims
	if plan.MaxClaims > 0 && record.ClaimCount >= plan.MaxClaims {
		return ClaimResultPayload{}, code.CheckInMaxReached.Err()
	}

	// 6. Compute streak for streak type
	newStreak := record.Streak
	if string(plan.Type) == "streak" {
		gracePeriod := time.Duration(plan.Interval*2) * time.Second
		if record.LastClaimed != nil && now.Sub(*record.LastClaimed) <= gracePeriod {
			newStreak++
		} else {
			newStreak = 1
		}
		// Reset streak if max_claims reached
		if plan.MaxClaims > 0 && newStreak > plan.MaxClaims {
			newStreak = 1
		}
	} else {
		newStreak = record.ClaimCount + 1
	}

	// 7. Determine reward
	reward := determineReward(plan, newStreak, record.ClaimCount+1)

	// 8. Ensure balance exists and grant reward
	s.service.AccountBalance.EnsureExists(ctx, accountId)
	if reward.Amount > 0 {
		s.service.repository.AccountBalance.AddCoins(ctx, accountId, reward.Amount)
	}

	// 9. Update record
	newClaimCount := record.ClaimCount + 1
	updatedRecord, err := s.service.repository.CheckInRecord.Update(ctx, checkInRecordRepo.UpdateParam{
		AccountId:   accountId,
		PlanId:      planId,
		ClaimCount:  newClaimCount,
		Streak:      newStreak,
		LastClaimed: now,
	})
	if err != nil {
		return ClaimResultPayload{}, code.HttpInternalServerError.Err().WithError(trace.Wrap(err))
	}

	// 10. Get balance
	balance, _ := s.service.AccountBalance.GetBalance(ctx, accountId)

	// 11. Compute status
	status := computeStatus(plan, updatedRecord)

	return ClaimResultPayload{
		Reward:  reward,
		Status:  status,
		Balance: balance,
	}, code.OK()
}

func computeStatus(plan model.CheckInPlan, record model.CheckInRecord) CheckInStatus {
	status := CheckInStatus{
		Streak:     record.Streak,
		ClaimCount: record.ClaimCount,
	}

	now := time.Now()

	if record.LastClaimed == nil {
		status.CanClaim = true
		return status
	}

	elapsed := now.Sub(*record.LastClaimed)
	if elapsed.Seconds() >= float64(plan.Interval) {
		status.CanClaim = true
	} else {
		status.CanClaim = false
		nextClaim := record.LastClaimed.Add(time.Duration(plan.Interval) * time.Second)
		status.NextClaimAt = &nextClaim
	}

	if plan.MaxClaims > 0 && record.ClaimCount >= plan.MaxClaims {
		status.CanClaim = false
		status.NextClaimAt = nil
	}

	return status
}

func determineReward(plan model.CheckInPlan, streak int32, claimNumber int32) CheckInReward {
	var rewards []CheckInReward
	_ = json.Unmarshal([]byte(plan.Rewards), &rewards)

	switch string(plan.Type) {
	case "recurring":
		if len(rewards) > 0 {
			return rewards[0]
		}
	case "streak":
		for _, r := range rewards {
			if r.Day == int(streak) {
				return r
			}
		}
		// Fallback to last reward
		if len(rewards) > 0 {
			return rewards[len(rewards)-1]
		}
	case "calendar":
		for _, r := range rewards {
			if r.Day == int(claimNumber) {
				return r
			}
		}
		// Fallback to a default coin reward
		return CheckInReward{Type: "coin", Amount: 10}
	}

	return CheckInReward{Type: "coin", Amount: 0}
}
