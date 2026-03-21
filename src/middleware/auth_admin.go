package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/database/zakiverse-db/public/model"
	"github.com/zakiverse/zakiverse-api/util/response"
)

func (m *Middleware) AuthAdmin(c *gin.Context) {
	role := c.GetString(cst.MiddlewareKeyAccountRole)
	if role != string(model.AccountRole_Admin) {
		response.Error(c, code.HttpForbidden, nil)
		c.Abort()
		return
	}
	c.Next()
}
