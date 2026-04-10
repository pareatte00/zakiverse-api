package card

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type findAllRequest struct {
	Search string `json:"search"`
	Rarity string `json:"rarity" validate:"omitempty,oneof=common rare epic legendary prismatic"`
	TagId  string `json:"tag_id" validate:"omitempty,uuid"`
	Sort   string `json:"sort" validate:"omitempty,oneof=name rarity"`
	Order  string `json:"order" validate:"omitempty,oneof=asc desc"`
	Page   int64  `json:"page" validate:"required,min=1"`
	Limit  int64  `json:"limit" validate:"required,min=1,max=100"`
}

func (h Handler) FindAll(c *gin.Context) {
	var request findAllRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.Card.FindAll(c.Request.Context(), service.FindAllCardsParam{
		Search: request.Search,
		Rarity: request.Rarity,
		TagId:  request.TagId,
		Sort:   request.Sort,
		Order:  request.Order,
		Page:   request.Page,
		Limit:  request.Limit,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
