package anime

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type createOneRequest struct {
	MalId      int32   `json:"mal_id" validate:"required"`
	Title      string  `json:"title" validate:"required"`
	Synopsis   *string `json:"synopsis"`
	CoverImage *string `json:"cover_image"`
}

func (h Handler) CreateOne(c *gin.Context) {
	var request createOneRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.Anime.CreateOne(c.Request.Context(), service.CreateAnimeParam{
		MalId:      request.MalId,
		Title:      request.Title,
		Synopsis:   request.Synopsis,
		CoverImage: request.CoverImage,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
