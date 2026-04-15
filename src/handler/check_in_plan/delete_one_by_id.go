package check_in_plan

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type deleteOneByIdUri struct {
	ID string `uri:"id" validate:"required,uuid"`
}

func (h Handler) DeleteOneById(c *gin.Context) {
	var uri deleteOneByIdUri
	if !binder.ShouldBindUri(c, &uri) {
		return
	}

	codeErr := h.service.CheckIn.DeletePlanById(c.Request.Context(), uri.ID)
	if !codeErr.OK() {
		response.Error(c, codeErr.Code(), response.NewError().WithDebug(codeErr.Error()))
		return
	}

	response.Http(c, http.StatusOK, nil)
}
