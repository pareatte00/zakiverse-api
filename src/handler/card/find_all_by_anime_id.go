package card

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type findAllByAnimeIdUri struct {
	AnimeId string `uri:"animeId" validate:"required,uuid"`
}

type findAllByAnimeIdQuery struct {
	Page  int64 `form:"page" validate:"required,min=1"`
	Limit int64 `form:"limit" validate:"required,min=1,max=100"`
}

func (h Handler) FindAllByAnimeId(c *gin.Context) {
	var uri findAllByAnimeIdUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	var query findAllByAnimeIdQuery
	if !binder.ShouldBindQuery(c, &query) {
		return
	}

	payload, codeErr := h.service.Card.FindAllByAnimeId(c.Request.Context(), service.FindAllCardsByAnimeIdParam{
		AnimeId: uri.AnimeId,
		Page:    query.Page,
		Limit:   query.Limit,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
