package pack

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type createOneRequest struct {
	Name         string             `json:"name" validate:"required"`
	Description  *string            `json:"description"`
	Image        string             `json:"image" validate:"required"`
	CardsPerPull int32              `json:"cards_per_pull" validate:"required,min=1,max=20"`
	IsActive     bool               `json:"is_active"`
	OpenAt       *time.Time         `json:"open_at"`
	CloseAt      *time.Time         `json:"close_at"`
	Config       service.PackConfig `json:"config" validate:"required"`
}

func (h Handler) CreateOne(c *gin.Context) {
	var request createOneRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.Pack.CreateOne(c.Request.Context(), service.CreatePackParam{
		Name:         request.Name,
		Description:  request.Description,
		Image:        request.Image,
		CardsPerPull: request.CardsPerPull,
		IsActive:     request.IsActive,
		OpenAt:       request.OpenAt,
		CloseAt:      request.CloseAt,
		Config:       request.Config,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
