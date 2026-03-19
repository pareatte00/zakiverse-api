package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type authDiscordRequest struct {
	Code string `json:"code" validate:"required"`
}

func (h Handler) AuthDiscord(c *gin.Context) {
	var request authDiscordRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.Account.AuthDiscord(c.Request.Context(), service.AuthDiscordParam{
		Code: request.Code,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
