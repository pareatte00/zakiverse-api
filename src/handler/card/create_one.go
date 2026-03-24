package card

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type createOneRequest struct {
	MalId   int32  `json:"mal_id" validate:"required"`
	AnimeId string `json:"anime_id" validate:"required,uuid"`
	Rarity  string `json:"rarity" validate:"required,oneof=common uncommon rare epic legendary"`
	Name    string `json:"name" validate:"required"`
	Image   string `json:"image" validate:"required"`
}

func (h Handler) CreateOne(c *gin.Context) {
	var request createOneRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.Card.CreateOne(c.Request.Context(), service.CreateCardParam{
		MalId:   request.MalId,
		AnimeId: request.AnimeId,
		Rarity:  request.Rarity,
		Name:    request.Name,
		Image:   request.Image,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
