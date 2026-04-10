package pack_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type findAllQuery struct {
	BannerType string `form:"banner_type" validate:"omitempty,oneof=standard featured event beginner seasonal"`
	ActiveOnly bool   `form:"active_only"`
	Page       int64  `form:"page" validate:"required,min=1"`
	Limit      int64  `form:"limit" validate:"required,min=1,max=100"`
}

func (h Handler) FindAll(c *gin.Context) {
	var query findAllQuery
	if !binder.ShouldBindQuery(c, &query) {
		return
	}

	payload, codeErr := h.service.PackPool.FindAll(c.Request.Context(), service.FindAllPackPoolsParam{
		BannerType: query.BannerType,
		ActiveOnly: query.ActiveOnly,
		Page:       query.Page,
		Limit:      query.Limit,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
