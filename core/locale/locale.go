package locale

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/cst"
)

type Locale string

const (
	EN Locale = "en"
	TH Locale = "th"
)

var supported = map[Locale]struct{}{
	EN: {},
	TH: {},
}

func Default() Locale {
	return EN
}

func Available(locale string) bool {
	_, ok := supported[Locale(strings.ToLower(locale))]
	return ok
}

func GetLocale(c *gin.Context) Locale {
	if v := c.GetString(cst.MiddlewareKeyLocale); v != "" && Available(v) {
		return Locale(v)
	}

	return EN
}
