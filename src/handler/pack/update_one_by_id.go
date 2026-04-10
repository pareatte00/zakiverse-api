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
	ID string `uri:"id" validate:"required,uuid"`
}

type updateOneByIdRequest struct {
	Code          *string             `json:"code" validate:"omitempty"`
	Name          *string             `json:"name" validate:"omitempty"`
	Description   *string             `json:"description"`
	Image         *string             `json:"image" validate:"omitempty"`
	NameImage     *string             `json:"name_image"`
	CardsPerPull  *int32              `json:"cards_per_pull" validate:"omitempty,min=1,max=20"`
	SortOrder     *int32              `json:"sort_order" validate:"omitempty"`
	Config        *service.PackConfig `json:"config"`
	RotationOrder *int32              `json:"rotation_order"`
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

	updates := patcher.Pick(updateMap, "code", "name", "description", "image", "name_image", "cards_per_pull", "sort_order", "config", "rotation_order")

	payload, codeErr := h.service.Pack.UpdateOneById(c.Request.Context(), uri.ID, updates)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
