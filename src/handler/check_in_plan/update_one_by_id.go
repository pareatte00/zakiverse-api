package check_in_plan

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/patcher"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type updateOneByIdUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

type updateOneByIdRequest struct {
	Code        *string    `json:"code"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Type        *string    `json:"type" validate:"omitempty,oneof=recurring streak calendar"`
	Interval    *int32     `json:"interval" validate:"omitempty,min=1"`
	MaxClaims   *int32     `json:"max_claims"`
	Rewards     *string    `json:"rewards"`
	ResetPolicy *string    `json:"reset_policy" validate:"omitempty,oneof=rolling daily_reset weekly_reset monthly_reset"`
	IsActive    *bool      `json:"is_active"`
	StartsAt    *time.Time `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
	SortOrder   *int32     `json:"sort_order"`
}

func (h Handler) UpdateOneById(c *gin.Context) {
	var uri updateOneByIdUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	var request updateOneByIdRequest
	if !binder.ShouldBindBufferedJson(c, &request) {
		return
	}

	var updateMap map[string]any
	if !binder.BindBufferedJson(c, &updateMap) {
		return
	}

	updates := patcher.Pick(updateMap, "code", "name", "description", "type", "interval", "max_claims", "rewards", "reset_policy", "is_active", "starts_at", "ends_at", "sort_order")

	payload, codeErr := h.service.CheckIn.UpdatePlanById(c.Request.Context(), uri.ID, updates)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
