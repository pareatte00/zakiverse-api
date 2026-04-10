package pack

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type removeCardsUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

type removeCardsRequest struct {
	CardIds []string `json:"card_ids" validate:"required,min=1,dive,uuid"`
}

func (h Handler) RemoveCards(c *gin.Context) {
	var uri removeCardsUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	var request removeCardsRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	codeErr := h.service.Pack.RemoveCards(c.Request.Context(), uri.ID, request.CardIds)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp())
}
