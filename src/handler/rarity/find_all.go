package rarity

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/response"
)

func (h Handler) FindAll(c *gin.Context) {
	payload, codeErr := h.service.Rarity.FindAll(c.Request.Context())
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
