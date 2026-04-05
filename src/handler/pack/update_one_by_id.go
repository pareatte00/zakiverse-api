package pack

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/patcher"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type updateOneByIdUri struct {
	Id string `uri:"id" validate:"required,uuid"`
}

type updateOneByIdRequest struct {
	Code         *string             `json:"code" validate:"omitempty"`
	Name         *string             `json:"name" validate:"omitempty"`
	Description  *string             `json:"description"`
	Image        *string             `json:"image" validate:"omitempty"`
	NameImage    *string             `json:"name_image"`
	Type         *string             `json:"type" validate:"omitempty,oneof=standard limited event"`
	CardsPerPull *int32              `json:"cards_per_pull" validate:"omitempty,min=1,max=20"`
	SortOrder    *int32              `json:"sort_order" validate:"omitempty"`
	IsActive     *bool               `json:"is_active"`
	OpenAt       *string             `json:"open_at"`
	CloseAt      *string             `json:"close_at"`
	Config       *service.PackConfig `json:"config"`
	PoolId       *string             `json:"pool_id" validate:"omitempty,uuid"`
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

	updates := patcher.Pick(updateMap, "code", "name", "description", "image", "name_image", "type", "cards_per_pull", "sort_order", "is_active", "open_at", "close_at", "config", "pool_id")

	payload, codeErr := h.service.Pack.UpdateOneById(c.Request.Context(), uri.Id, updates)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
