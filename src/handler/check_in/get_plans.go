package check_in

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/util/response"
)

func (h Handler) GetPlans(c *gin.Context) {
	accountId := c.GetString(cst.MiddlewareKeyAccountId)

	payload, codeErr := h.service.CheckIn.GetActivePlans(c.Request.Context(), accountId)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
