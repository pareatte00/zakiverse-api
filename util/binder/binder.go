package binder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/logger"
	"github.com/zakiverse/zakiverse-api/util/response"
	"github.com/zakiverse/zakiverse-api/util/validator"
)

func validate(c *gin.Context, obj any, httpStatus int) bool {
	isValidatePass, errorFieldList, validatorErr := validator.Validate(obj)
	if validatorErr != nil {
		logger.Error(cst.KeyValidator, logger.Field(cst.KeyError, validatorErr.Error()))
		response.Http(c, http.StatusInternalServerError, response.NewHttp().WithDebug(validatorErr))
		return false
	}
	if !isValidatePass {
		if httpStatus == http.StatusUnauthorized {
			response.Http(c, httpStatus, nil)
		} else {
			response.Http(c, httpStatus, response.NewHttp().WithMeta(errorFieldList))
		}
		return false
	}

	return true
}

func BindHeader(c *gin.Context, obj any) bool {
	if err := c.ShouldBindHeader(obj); err != nil {
		response.Http(c, http.StatusUnauthorized, nil)
		return false
	}

	return true
}

func ShouldBindHeader(c *gin.Context, obj any) bool {
	if !BindHeader(c, obj) {
		return false
	}
	if !validate(c, obj, http.StatusUnauthorized) {
		return false
	}

	return true
}

func BindBufferedJson(c *gin.Context, obj any) bool {
	if err := c.ShouldBindBodyWith(obj, binding.JSON); err != nil {
		response.Http(c, http.StatusBadRequest, response.NewHttp().WithMeta(err.Error()))
		return false
	}

	return true
}

func ShouldBindBufferedJson(c *gin.Context, obj any) bool {
	if !BindBufferedJson(c, obj) {
		return false
	}
	if !validate(c, obj, http.StatusBadRequest) {
		return false
	}

	return true
}

func BindJson(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.Http(c, http.StatusBadRequest, response.NewHttp().WithMeta(err.Error()))
		return false
	}

	return true
}

func ShouldBindJson(c *gin.Context, obj any) bool {
	if !BindJson(c, obj) {
		return false
	}
	if !validate(c, obj, http.StatusBadRequest) {
		return false
	}

	return true
}

func BindUri(c *gin.Context, obj any) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		response.Http(c, http.StatusBadRequest, response.NewHttp().WithMeta(err.Error()))
		return false
	}

	return true
}

func ShouldBindUri(c *gin.Context, obj any) bool {
	if !BindUri(c, obj) {
		return false
	}
	if !validate(c, obj, http.StatusBadRequest) {
		return false
	}

	return true
}

func BindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		response.Http(c, http.StatusBadRequest, response.NewHttp().WithMeta(err.Error()))
		return false
	}

	return true
}

func ShouldBindQuery(c *gin.Context, obj any) bool {
	if !BindQuery(c, obj) {
		return false
	}
	if !validate(c, obj, http.StatusBadRequest) {
		return false
	}

	return true
}
