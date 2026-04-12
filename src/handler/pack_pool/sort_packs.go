package pack_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type sortPacksUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

type sortPacksRequest struct {
	Ids []string `json:"ids" validate:"required,min=1,dive,uuid"`
}

func (h Handler) SortPacks(c *gin.Context) {
	var uri sortPacksUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	var request sortPacksRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	codeErr := h.service.PackPool.ReorderPacks(c.Request.Context(), uri.ID, request.Ids)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp())
}
