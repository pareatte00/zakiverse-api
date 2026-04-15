package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type getByIdOrNameUri struct {
	Identifier string `uri:"identifier" validate:"required"`
}

func (h Handler) GetByIdOrName(c *gin.Context) {
	var uri getByIdOrNameUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	payload, codeErr := h.service.Profile.GetProfile(c.Request.Context(), uri.Identifier)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
