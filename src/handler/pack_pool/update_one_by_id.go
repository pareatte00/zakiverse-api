package pack_pool

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
	Name              *string    `json:"name" validate:"omitempty"`
	Description       *string    `json:"description"`
	Image             *string    `json:"image"`
	BannerType        *string    `json:"banner_type" validate:"omitempty,oneof=standard featured"`
	SortOrder         *int32     `json:"sort_order"`
	IsActive          *bool      `json:"is_active"`
	OpenAt            *time.Time `json:"open_at"`
	CloseAt           *time.Time `json:"close_at"`
	ActiveCount       *int32     `json:"active_count" validate:"omitempty,min=1"`
	RotationType      *string    `json:"rotation_type" validate:"omitempty,oneof=none weekly monthly"`
	RotationDay       *int32     `json:"rotation_day" validate:"omitempty,min=0,max=31"`
	RotationInterval  *int32     `json:"rotation_interval" validate:"omitempty,min=1"`
	RotationHour      *int32     `json:"rotation_hour" validate:"omitempty,min=0,max=23"`
	TimezoneOffset    *int32     `json:"timezone_offset" validate:"omitempty,min=-12,max=14"`
	RotationOrderMode *string    `json:"rotation_order_mode" validate:"omitempty,oneof=auto manual"`
	PreviewDays       *int32     `json:"preview_days" validate:"omitempty,min=0"`
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

	updates := patcher.Pick(updateMap,
		"name", "description", "image", "banner_type", "sort_order",
		"is_active", "open_at", "close_at", "active_count",
		"rotation_type", "rotation_day", "rotation_interval", "rotation_hour",
		"timezone_offset", "rotation_order_mode", "preview_days",
	)

	payload, codeErr := h.service.PackPool.UpdateOneById(c.Request.Context(), uri.ID, updates)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
