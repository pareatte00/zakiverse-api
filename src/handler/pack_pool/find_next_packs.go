package pack_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type findNextPacksUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

func (h Handler) FindNextPacks(c *gin.Context) {
	var uri findNextPacksUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	payload, codeErr := h.service.PackPool.FindNextPacks(c.Request.Context(), uri.ID)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
