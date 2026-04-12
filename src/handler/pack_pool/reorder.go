package pack_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type reorderRequest struct {
	BannerType string   `json:"banner_type" validate:"required,oneof=standard featured event beginner seasonal"`
	Ids        []string `json:"ids" validate:"required,min=1,dive,uuid"`
}

func (h Handler) Reorder(c *gin.Context) {
	var request reorderRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	codeErr := h.service.PackPool.Reorder(c.Request.Context(), request.BannerType, request.Ids)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp())
}
