package check_in_plan

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type createOneRequest struct {
	Code        string     `json:"code" validate:"required"`
	Name        string     `json:"name" validate:"required"`
	Description *string    `json:"description"`
	Type        string     `json:"type" validate:"required,oneof=recurring streak calendar"`
	Interval    int32      `json:"interval" validate:"required,min=1"`
	MaxClaims   int32      `json:"max_claims"`
	Rewards     string     `json:"rewards" validate:"required"`
	ResetPolicy string     `json:"reset_policy" validate:"required,oneof=rolling daily_reset weekly_reset monthly_reset"`
	IsActive    bool       `json:"is_active"`
	StartsAt    *time.Time `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
	SortOrder   int32      `json:"sort_order"`
}

func (h Handler) CreateOne(c *gin.Context) {
	var request createOneRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.CheckIn.CreatePlan(c.Request.Context(), service.CreateCheckInPlanParam{
		Code:        request.Code,
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Interval:    request.Interval,
		MaxClaims:   request.MaxClaims,
		Rewards:     request.Rewards,
		ResetPolicy: request.ResetPolicy,
		IsActive:    request.IsActive,
		StartsAt:    request.StartsAt,
		EndsAt:      request.EndsAt,
		SortOrder:   request.SortOrder,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
