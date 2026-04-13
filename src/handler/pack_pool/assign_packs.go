package pack_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type assignPacksUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

type assignPacksRequest struct {
	Ids []string `json:"ids" validate:"dive,uuid"`
}

func (h Handler) AssignPacks(c *gin.Context) {
	var uri assignPacksUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	var request assignPacksRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	codeErr := h.service.PackPool.AssignPacks(c.Request.Context(), uri.ID, request.Ids)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp())
}
