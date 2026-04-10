package pack

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type addCardsUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

type addCardsRequestItem struct {
	CardId       string   `json:"card_id" validate:"required,uuid"`
	Weight       float64  `json:"weight"`
	IsFeatured   bool     `json:"is_featured"`
	FeaturedRate *float64 `json:"featured_rate" validate:"omitempty,min=0,max=1"`
}

type addCardsRequest struct {
	Cards []addCardsRequestItem `json:"cards" validate:"required,min=1,dive"`
}

func (h Handler) AddCards(c *gin.Context) {
	var uri addCardsUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	var request addCardsRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	params := make([]service.AddPackCardsParam, len(request.Cards))
	for i, card := range request.Cards {
		params[i] = service.AddPackCardsParam{
			CardId:       card.CardId,
			Weight:       card.Weight,
			IsFeatured:   card.IsFeatured,
			FeaturedRate: card.FeaturedRate,
		}
	}

	payload, codeErr := h.service.Pack.AddCards(c.Request.Context(), uri.ID, params)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
