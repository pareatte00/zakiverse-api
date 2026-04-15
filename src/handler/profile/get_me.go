package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/util/response"
)

func (h Handler) GetMe(c *gin.Context) {
	accountId := c.GetString(cst.MiddlewareKeyAccountId)

	payload, codeErr := h.service.Profile.GetProfile(c.Request.Context(), accountId)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
