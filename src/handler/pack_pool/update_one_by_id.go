package pack_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/patcher"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type updateOneByIdUri struct {
	Id string `uri:"id" validate:"required,uuid"`
}

type updateOneByIdRequest struct {
	Name        *string `json:"name" validate:"omitempty"`
	Description *string `json:"description"`
	ActiveCount *int32  `json:"active_count" validate:"omitempty,min=1"`
	RotationDay *int32  `json:"rotation_day" validate:"omitempty,min=0,max=6"`
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

	updates := patcher.Pick(updateMap, "name", "description", "active_count", "rotation_day")

	payload, codeErr := h.service.PackPool.UpdateOneById(c.Request.Context(), uri.Id, updates)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
