package code

//go:generate go run github.com/zakiverse/zakiverse-api/cmd/codegen/code

import (
	"net/http"

	"github.com/zakiverse/zakiverse-api/core/locale"
)

type Code string

type Definition struct {
	StatusCode int
	Message    map[locale.Locale]string
}

type I struct {
	code Code
	err  error
}

func (c I) OK() bool {
	return c.code == ""
}

func New(code Code) I {
	return I{
		code: code,
	}
}

func (c I) WithError(err error) I {
	return I{
		code: c.code,
		err:  err,
	}
}

func (c I) Code() Code {
	return c.code
}

func (c I) Error() error {
	return c.err
}

func (c I) ErrorText() string {
	if c.err != nil {
		return c.err.Error()
	}

	return ""
}

var registry = map[Code]Definition{}

func reg(code Code, def Definition) {
	if _, exists := registry[code]; exists {
		panic("duplicate error code: " + string(code))
	}
	if def.StatusCode == 0 {
		panic("status code is unset (0): " + string(code))
	}

	registry[code] = def
}

func GetStatusCode(code Code) (Code, int) {
	if def, ok := registry[code]; ok {
		return code, def.StatusCode
	}

	return ErrorNotImplemented, http.StatusNotImplemented
}

func GetMessage(code Code, loc locale.Locale) string {
	def, ok := registry[code]
	if !ok {
		def = registry[ErrorNotImplemented]
	}
	if msg, ok := def.Message[loc]; ok {
		return msg
	}

	return registry[ErrorNotImplemented].Message[locale.EN]
}
