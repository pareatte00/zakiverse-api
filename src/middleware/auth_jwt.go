package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/core/cst"
	jwtutil "github.com/zakiverse/zakiverse-api/util/jwt"
	"github.com/zakiverse/zakiverse-api/util/response"
)

func (m *Middleware) AuthJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, code.HttpUnauthorized, nil)
		c.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.Error(c, code.HttpUnauthorized, nil)
		c.Abort()
		return
	}

	claims, err := jwtutil.Parse(parts[1], m.credential.JwtSecret)
	if err != nil {
		response.Error(c, code.HttpUnauthorized, nil)
		c.Abort()
		return
	}

	c.Set(cst.MiddlewareKeyAccountId, claims.AccountId.String())
	c.Set(cst.MiddlewareKeyAccountRole, claims.Role)
	c.Next()
}
