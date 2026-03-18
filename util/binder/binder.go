package binder

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/logger"
	"github.com/zakiverse/zakiverse-api/util/response"
	"github.com/zakiverse/zakiverse-api/util/validator"
)

func validate(c *gin.Context, obj any, co code.Code) bool {
	isValidatePass, errorFieldList, validatorErr := validator.Validate(obj)
	if validatorErr != nil {
		logger.Error(cst.KeyValidator, logger.Field(cst.KeyError, validatorErr.Error()))
		response.Json(c, code.HttpInternalServerError, response.NewParam().
			WithDebug(validatorErr),
		)
		return false
	}
	if !isValidatePass {
		if co == code.HttpUnauthorized {
			response.Json(c, co, nil)
		} else {
			response.Json(c, co, response.NewParam().
				WithMeta(errorFieldList),
			)
		}
		return false
	}

	return true
}

func BindHeader(c *gin.Context, obj any) bool {
	if err := c.ShouldBindHeader(obj); err != nil {
		response.Json(c, code.HttpUnauthorized, nil)
		return false
	}

	return true
}

func ShouldBindHeader(c *gin.Context, obj any) bool {
	if !BindHeader(c, obj) {
		return false
	}
	if !validate(c, obj, code.HttpUnauthorized) {
		return false
	}

	return true
}

func BindBufferedJson(c *gin.Context, obj any) bool {
	if err := c.ShouldBindBodyWith(obj, binding.JSON); err != nil {
		response.Json(c, code.HttpBadRequest, response.NewParam().
			WithMeta(err.Error()),
		)
		return false
	}

	return true
}

func ShouldBindBufferedJson(c *gin.Context, obj any) bool {
	if !BindBufferedJson(c, obj) {
		return false
	}
	if !validate(c, obj, code.HttpBadRequest) {
		return false
	}

	return true
}

func BindJson(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.Json(c, code.HttpBadRequest, response.NewParam().
			WithMeta(err.Error()),
		)
		return false
	}

	return true
}

func ShouldBindJson(c *gin.Context, obj any) bool {
	if !BindJson(c, obj) {
		return false
	}
	if !validate(c, obj, code.HttpBadRequest) {
		return false
	}

	return true
}

func BindUri(c *gin.Context, obj any) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		response.Json(c, code.HttpBadRequest, response.NewParam().
			WithMeta(err.Error()),
		)
		return false
	}

	return true
}

func ShouldBindUri(c *gin.Context, obj any) bool {
	if !BindUri(c, obj) {
		return false
	}
	if !validate(c, obj, code.HttpBadRequest) {
		return false
	}

	return true
}

func BindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		response.Json(c, code.HttpBadRequest, response.NewParam().
			WithMeta(err.Error()),
		)
		return false
	}

	return true
}

func ShouldBindQuery(c *gin.Context, obj any) bool {
	if !BindQuery(c, obj) {
		return false
	}
	if !validate(c, obj, code.HttpBadRequest) {
		return false
	}

	return true
}
