package response

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zakiverse/zakiverse-api/core/code"
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/core/locale"
)

type Response struct {
	Timestamp string `json:"timestamp"`
	Detail    Detail `json:"detail"`
	Payload   any    `json:"payload,omitempty"`
	Debug     any    `json:"debug,omitempty"`
	Version   string `json:"version"`
}

type Param struct {
	payload any
	debug   error
	meta    any
}

type Detail struct {
	Code    code.Code `json:"code"`
	Message string    `json:"message"`
	Meta    any       `json:"meta,omitempty"`
}

func NewParam() *Param {
	return &Param{}
}

func (p *Param) WithPayload(payload any) *Param {
	p.payload = payload
	return p
}

func (p *Param) WithDebug(debug error) *Param {
	p.debug = debug
	return p
}

func (p *Param) WithMeta(meta any) *Param {
	p.meta = meta
	return p
}

func Json(c *gin.Context, co code.Code, param *Param) {
	loc := locale.GetLocale(c)
	cod, statusCode := code.GetStatusCode(co)
	msg := code.GetMessage(co, loc)

	var payload any
	if param != nil && param.payload != nil {
		payload = param.payload
	}

	var debug any
	if strings.ToLower(viper.GetString("application.deploy_mode")) == cst.DeployModeDevelopment {
		if param != nil && param.debug != nil {
			debug = param.debug.Error()
		}
	}

	var meta any
	if param != nil && param.meta != nil {
		meta = param.meta
	}

	c.JSON(statusCode, Response{
		Timestamp: time.Now().Format(time.RFC3339),
		Payload:   payload,
		Detail: Detail{
			Code:    cod,
			Message: msg,
			Meta:    meta,
		},
		Debug:   debug,
		Version: os.Getenv("APP_VERSION"),
	})
}
