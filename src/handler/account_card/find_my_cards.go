package account_card

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type findMyCardsQuery struct {
	Page  int64 `form:"page" validate:"required,min=1"`
	Limit int64 `form:"limit" validate:"required,min=1,max=100"`
}

func (h Handler) FindMyCards(c *gin.Context) {
	var query findMyCardsQuery
	if !binder.ShouldBindQuery(c, &query) {
		return
	}

	accountId := c.GetString(cst.MiddlewareKeyAccountId)

	payload, meta, codeErr := h.service.AccountCard.FindMyCards(c.Request.Context(), service.FindMyCardsParam{
		AccountId: accountId,
		Page:      query.Page,
		Limit:     query.Limit,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload).WithMeta(meta))
}
