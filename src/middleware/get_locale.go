package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/core/locale"
	"github.com/zakiverse/zakiverse-api/util/binder"
)

type getLocaleHeader struct {
	Locale string `header:"X-Locale"`
}

func (m *Middleware) GetLocale(c *gin.Context) {
	var header getLocaleHeader

	if !binder.BindHeader(c, &header) {
		c.Abort()
		return
	}

	loc := strings.ToLower(header.Locale)

	if locale.Available(loc) {
		c.Set(cst.MiddlewareKeyLocale, loc)
	} else {
		c.Set(cst.MiddlewareKeyLocale, locale.EN)
	}

	c.Next()
}
