package pack_pool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type createOneRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	ActiveCount int32   `json:"active_count" validate:"required,min=1"`
	RotationDay int32   `json:"rotation_day" validate:"required,min=0,max=6"`
}

func (h Handler) CreateOne(c *gin.Context) {
	var request createOneRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.PackPool.CreateOne(c.Request.Context(), service.CreatePackPoolParam{
		Name:        request.Name,
		Description: request.Description,
		ActiveCount: request.ActiveCount,
		RotationDay: request.RotationDay,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
