package anime

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type updateOneByIdUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

type updateOneByIdRequest struct {
	Title      string  `json:"title" validate:"required"`
	Synopsis   *string `json:"synopsis"`
	CoverImage *string `json:"cover_image"`
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

	payload, codeErr := h.service.Anime.UpdateOneById(c.Request.Context(), uri.ID, service.UpdateAnimeParam{
		Title:      request.Title,
		Synopsis:   request.Synopsis,
		CoverImage: request.CoverImage,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
