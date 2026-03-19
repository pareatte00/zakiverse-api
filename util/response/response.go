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

type HttpResponse struct {
	Timestamp string `json:"timestamp"`
	Payload   any    `json:"payload,omitempty"`
	Meta      any    `json:"meta,omitempty"`
	Debug     any    `json:"debug,omitempty"`
	Version   string `json:"version"`
}

type ErrorResponse struct {
	Timestamp string      `json:"timestamp"`
	Error     ErrorDetail `json:"error"`
	Meta      any         `json:"meta,omitempty"`
	Debug     any         `json:"debug,omitempty"`
	Version   string      `json:"version"`
}

type ErrorDetail struct {
	Code    code.Code `json:"code"`
	Message string    `json:"message"`
}

type HttpParam struct {
	payload any
	meta    any
	debug   error
}

type ErrorParam struct {
	meta  any
	debug error
}

func NewHttp() *HttpParam {
	return &HttpParam{}
}

func (p *HttpParam) WithPayload(payload any) *HttpParam {
	p.payload = payload
	return p
}

func (p *HttpParam) WithDebug(debug error) *HttpParam {
	p.debug = debug
	return p
}

func (p *HttpParam) WithMeta(meta any) *HttpParam {
	p.meta = meta
	return p
}

func NewError() *ErrorParam {
	return &ErrorParam{}
}

func (p *ErrorParam) WithDebug(debug error) *ErrorParam {
	p.debug = debug
	return p
}

func (p *ErrorParam) WithMeta(meta any) *ErrorParam {
	p.meta = meta
	return p
}

func Http(c *gin.Context, statusCode int, param *HttpParam) {
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

	c.JSON(statusCode, HttpResponse{
		Timestamp: time.Now().Format(time.RFC3339),
		Payload:   payload,
		Meta:      meta,
		Debug:     debug,
		Version:   os.Getenv("APP_VERSION"),
	})
}

func Error(c *gin.Context, co code.Code, param *ErrorParam) {
	loc := locale.GetLocale(c)
	cod, statusCode := code.GetStatusCode(co)
	msg := code.GetMessage(co, loc)

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

	c.JSON(statusCode, ErrorResponse{
		Timestamp: time.Now().Format(time.RFC3339),
		Error: ErrorDetail{
			Code:    cod,
			Message: msg,
		},
		Meta:    meta,
		Debug:   debug,
		Version: os.Getenv("APP_VERSION"),
	})
}
