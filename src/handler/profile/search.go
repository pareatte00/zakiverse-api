package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type searchRequest struct {
	Query string `json:"query" validate:"required,min=1"`
	Page  int64  `json:"page" validate:"required,min=1"`
	Limit int64  `json:"limit" validate:"required,min=1,max=100"`
}

func (h Handler) Search(c *gin.Context) {
	var request searchRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, meta, codeErr := h.service.Profile.SearchProfiles(c.Request.Context(), request.Query, request.Page, request.Limit)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload).WithMeta(meta))
}
