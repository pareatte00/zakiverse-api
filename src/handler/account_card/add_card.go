package account_card

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type addCardRequest struct {
	AccountId string `json:"account_id" validate:"required,uuid"`
	CardId    string `json:"card_id" validate:"required,uuid"`
}

func (h Handler) AddCard(c *gin.Context) {
	var request addCardRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.AccountCard.AddCard(c.Request.Context(), service.AddCardParam{
		AccountId: request.AccountId,
		CardId:    request.CardId,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
