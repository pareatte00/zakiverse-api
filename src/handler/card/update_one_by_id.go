package card

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type updateOneByIdUri struct {
	Id string `uri:"id" validate:"required,uuid"`
}

type updateOneByIdRequest struct {
	Rarity string `json:"rarity" validate:"required,oneof=common uncommon rare epic legendary"`
	Name   string `json:"name" validate:"required"`
	Image  string `json:"image" validate:"required"`
}

func (h Handler) UpdateOneById(c *gin.Context) {
	var uri updateOneByIdUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	var request updateOneByIdRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.Card.UpdateOneById(c.Request.Context(), uri.Id, service.UpdateCardParam{
		Rarity: request.Rarity,
		Name:   request.Name,
		Image:  request.Image,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
