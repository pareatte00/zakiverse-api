package pack

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cst "github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type getPullHistoryUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

type getPullHistoryQuery struct {
	Page  int64 `form:"page" validate:"required,min=1"`
	Limit int64 `form:"limit" validate:"required,min=1,max=100"`
}

func (h Handler) GetPullHistory(c *gin.Context) {
	var uri getPullHistoryUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	var query getPullHistoryQuery
	if !binder.ShouldBindQuery(c, &query) {
		return
	}

	accountId := c.GetString(cst.MiddlewareKeyAccountId)

	payload, codeErr := h.service.Pack.GetPullHistory(c.Request.Context(), service.FindPullHistoryParam{
		AccountId: accountId,
		PackId:    uri.ID,
		Page:      query.Page,
		Limit:     query.Limit,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
