package anime

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type deleteOneByIdUri struct {
	Id string `uri:"id" validate:"required,uuid"`
}

func (h Handler) DeleteOneById(c *gin.Context) {
	var uri deleteOneByIdUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	codeErr := h.service.Anime.DeleteOneById(c.Request.Context(), uri.Id)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, nil)
}
