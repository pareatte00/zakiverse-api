package pack_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type findOneWithPacksUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

func (h Handler) FindOneWithPacks(c *gin.Context) {
	var uri findOneWithPacksUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	payload, codeErr := h.service.PackPool.FindOneByIdWithPacks(c.Request.Context(), uri.ID)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
