package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zakiverse/zakiverse-api/util/response"
)

func empty(c *gin.Context) {
	response.Http(c, http.StatusNotFound, nil)
}

func health(c *gin.Context) {
	response.Http(c, http.StatusOK, nil)
}

type InfoPayload struct {
	Version string `json:"version"`
}

func info(c *gin.Context) {
	response.Http(c, http.StatusOK, response.NewHttp().WithPayload(InfoPayload{
		Version: os.Getenv("APP_VERSION"),
	}),
	)
}
