package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/src/service"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type updateMeRequest struct {
	DisplayName   string   `json:"display_name" validate:"required,min=3,max=50"`
	Bio           *string  `json:"bio" validate:"omitempty,max=500"`
	ShowcaseCards []string `json:"showcase_cards" validate:"omitempty,dive,uuid"`
}

func (h Handler) UpdateMe(c *gin.Context) {
	var request updateMeRequest
	if !binder.ShouldBindJson(c, &request) {
		return
	}

	accountId := c.GetString(cst.MiddlewareKeyAccountId)

	payload, codeErr := h.service.Profile.UpdateProfile(c.Request.Context(), service.UpdateProfileParam{
		AccountId:     accountId,
		DisplayName:   request.DisplayName,
		Bio:           request.Bio,
		ShowcaseCards: request.ShowcaseCards,
	})
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(payload))
}
