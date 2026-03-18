package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/util/response"
)

func empty(c *gin.Context) {
	response.Json(c, code.HttpNotFound, nil)
}

func health(c *gin.Context) {
	response.Json(c, code.HttpOK, nil)
}

type InfoPayload struct {
	Version string `json:"version"`
}

func info(c *gin.Context) {
	response.Json(c, code.HttpOK, response.NewParam().
		WithPayload(InfoPayload{
			Version: os.Getenv("APP_VERSION"),
		}),
	)
}
