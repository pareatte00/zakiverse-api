package pack_pool

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type createOneRequest struct {
	Name              string     `json:"name" validate:"required"`
	Description       *string    `json:"description"`
	Image             *string    `json:"image"`
	BannerType        string     `json:"banner_type" validate:"required,oneof=standard featured event beginner seasonal"`
	SortOrder         int32      `json:"sort_order"`
	IsActive          bool       `json:"is_active"`
	OpenAt            *time.Time `json:"open_at"`
	CloseAt           *time.Time `json:"close_at"`
	ActiveCount       int32      `json:"active_count" validate:"required,min=1"`
	RotationType      string     `json:"rotation_type" validate:"required,oneof=none weekly monthly"`
	RotationDay       *int32     `json:"rotation_day" validate:"omitempty,min=0,max=31"`
	RotationInterval  int32      `json:"rotation_interval" validate:"min=1"`
	RotationHour      int32      `json:"rotation_hour" validate:"min=0,max=23"`
	RotationOrderMode string     `json:"rotation_order_mode" validate:"required,oneof=auto manual"`
	PreviewDays       int32      `json:"preview_days" validate:"min=0"`
}

func (h Handler) CreateOne(c *gin.Context) {
	var request createOneRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.PackPool.CreateOne(c.Request.Context(), service.CreatePackPoolParam{
		Name:              request.Name,
		Description:       request.Description,
		Image:             request.Image,
		BannerType:        request.BannerType,
		SortOrder:         request.SortOrder,
		IsActive:          request.IsActive,
		OpenAt:            request.OpenAt,
		CloseAt:           request.CloseAt,
		ActiveCount:       request.ActiveCount,
		RotationType:      request.RotationType,
		RotationDay:       request.RotationDay,
		RotationInterval:  request.RotationInterval,
		RotationHour:      request.RotationHour,
		RotationOrderMode: request.RotationOrderMode,
		PreviewDays:       request.PreviewDays,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
