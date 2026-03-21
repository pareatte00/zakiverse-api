package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/response"
)

func (h Handler) FindOneById(c *gin.Context) {
	accountId, _ := uuid.Parse(c.GetString(cst.MiddlewareKeyAccountId))

	payload, codeErr := h.service.Account.FindOneById(c.Request.Context(), service.FindOneByIdParam{
		AccountId: accountId,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
