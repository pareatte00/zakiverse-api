package card_tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type findAllRequest struct {
	Search string `json:"search"`
	Page   int64  `json:"page" validate:"required,min=1"`
	Limit  int64  `json:"limit" validate:"required,min=1,max=100"`
}

func (h Handler) FindAll(c *gin.Context) {
	var request findAllRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, meta, codeErr := h.service.CardTag.FindAll(c.Request.Context(), service.FindAllCardTagsParam{
		Search: request.Search,
		Page:   request.Page,
		Limit:  request.Limit,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload).WithMeta(meta))
}
