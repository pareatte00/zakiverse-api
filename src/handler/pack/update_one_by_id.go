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
	Name         *string             `json:"name" validate:"omitempty"`
	Description  *string             `json:"description"`
	Image        *string             `json:"image" validate:"omitempty"`
	CardsPerPull *int32              `json:"cards_per_pull" validate:"omitempty,min=1,max=20"`
	IsActive     *bool               `json:"is_active"`
	OpenAt       *string             `json:"open_at"`
	CloseAt      *string             `json:"close_at"`
	Config       *service.PackConfig `json:"config"`
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

	updates := patcher.Pick(updateMap, "name", "description", "image", "cards_per_pull", "is_active", "open_at", "close_at", "config")

	payload, codeErr := h.service.Pack.UpdateOneById(c.Request.Context(), uri.Id, updates)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
