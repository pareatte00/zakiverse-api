package pack

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cst "github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type pullUri struct {
	Id string `uri:"id" validate:"required,uuid"`
}

func (h Handler) Pull(c *gin.Context) {
	var uri pullUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	accountId := c.GetString(cst.MiddlewareKeyAccountId)

	payload, codeErr := h.service.Pack.Pull(c.Request.Context(), accountId, uri.Id)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
