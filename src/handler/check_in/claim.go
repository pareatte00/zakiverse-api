package check_in

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type claimUri struct {
	PlanId string `uri:"planId" validate:"required,uuid"`
}

func (h Handler) Claim(c *gin.Context) {
	var uri claimUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	accountId := c.GetString(cst.MiddlewareKeyAccountId)

	payload, codeErr := h.service.CheckIn.Claim(c.Request.Context(), accountId, uri.PlanId)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
