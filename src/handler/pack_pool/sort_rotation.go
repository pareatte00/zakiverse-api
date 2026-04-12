package pack_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type sortRotationUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

type sortRotationRequest struct {
	Ids []string `json:"ids" validate:"required,min=1,dive,uuid"`
}

func (h Handler) SortRotation(c *gin.Context) {
	var uri sortRotationUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	var request sortRotationRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	codeErr := h.service.PackPool.ReorderRotation(c.Request.Context(), uri.ID, request.Ids)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp())
}
