package card

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type createOneRequest struct {
	MalId           int32              `json:"mal_id" validate:"required"`
	Rarity          string             `json:"rarity" validate:"required,oneof=common rare epic legendary prismatic"`
	Name            string             `json:"name" validate:"required"`
	Image           string             `json:"image" validate:"required"`
	Config          service.CardConfig `json:"config"`
	AnimeMalId      int32              `json:"anime_mal_id" validate:"required"`
	AnimeTitle      string             `json:"anime_title" validate:"required"`
	AnimeSynopsis   *string            `json:"anime_synopsis"`
	AnimeCoverImage *string            `json:"anime_cover_image"`
}

func (h Handler) CreateOne(c *gin.Context) {
	var request createOneRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.Card.CreateOne(c.Request.Context(), service.CreateCardParam{
		MalId:           request.MalId,
		Rarity:          request.Rarity,
		Name:            request.Name,
		Image:           request.Image,
		Config:          request.Config,
		AnimeMalId:      request.AnimeMalId,
		AnimeTitle:      request.AnimeTitle,
		AnimeSynopsis:   request.AnimeSynopsis,
		AnimeCoverImage: request.AnimeCoverImage,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
