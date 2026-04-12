package card_tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type findAllQuery struct {
	Page  int64 `form:"page" validate:"required,min=1"`
	Limit int64 `form:"limit" validate:"required,min=1,max=100"`
}

func (h Handler) FindAll(c *gin.Context) {
	var query findAllQuery
	if !binder.ShouldBindQuery(c, &query) {
		return
	}

	payload, meta, codeErr := h.service.CardTag.FindAll(c.Request.Context(), service.FindAllCardTagsParam{
		Page:  query.Page,
		Limit: query.Limit,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload).WithMeta(meta))
}
