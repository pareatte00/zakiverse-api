package pack

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type createOneRequest struct {
	Code          string             `json:"code" validate:"required"`
	Name          string             `json:"name" validate:"required"`
	Description   *string            `json:"description"`
	Image         string             `json:"image" validate:"required"`
	NameImage     *string            `json:"name_image"`
	CardsPerPull  int32              `json:"cards_per_pull" validate:"required,min=1,max=20"`
	SortOrder     int32              `json:"sort_order"`
	Config        service.PackConfig `json:"config" validate:"required"`
	PoolId        string             `json:"pool_id" validate:"required,uuid"`
	RotationOrder *int32             `json:"rotation_order"`
}

func (h Handler) CreateOne(c *gin.Context) {
	var request createOneRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	payload, codeErr := h.service.Pack.CreateOne(c.Request.Context(), service.CreatePackParam{
		Code:          request.Code,
		Name:          request.Name,
		Description:   request.Description,
		Image:         request.Image,
		NameImage:     request.NameImage,
		CardsPerPull:  request.CardsPerPull,
		SortOrder:     request.SortOrder,
		Config:        request.Config,
		PoolId:        request.PoolId,
		RotationOrder: request.RotationOrder,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusCreated, response.NewHttp().WithPayload(payload))
}
