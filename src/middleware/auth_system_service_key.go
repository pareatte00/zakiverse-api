package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/util/binder"
	"github.com/zakiverse/zakiverse-api/util/response"
)

type authSystemServiceKeyHeader struct {
	SystemServiceKey string `header:"X-System-Key" validate:"required"`
}

func (m *Middleware) AuthSystemServiceKey(c *gin.Context) {
	var header authSystemServiceKeyHeader

	if !binder.ShouldBindHeader(c, &header) {
		c.Abort()
		return
	}

	if m.credential.SystemServiceKey == "" || m.credential.SystemServiceKey != header.SystemServiceKey {
		response.Json(c, code.HttpUnauthorized, nil)
		c.Abort()
		return
	}

	c.Next()
}
