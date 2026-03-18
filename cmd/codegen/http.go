package codegen

import (
	"fmt"
	"strconv"
)

var ValidHttpStatusNames = map[string]bool{
	"StatusContinue":              true,
	"StatusSwitchingProtocols":    true,
	"StatusOK":                    true,
	"StatusCreated":               true,
	"StatusAccepted":              true,
	"StatusNoContent":             true,
	"StatusMovedPermanently":      true,
	"StatusFound":                 true,
	"StatusNotModified":           true,
	"StatusTemporaryRedirect":     true,
	"StatusPermanentRedirect":     true,
	"StatusBadRequest":            true,
	"StatusUnauthorized":          true,
	"StatusForbidden":             true,
	"StatusNotFound":              true,
	"StatusMethodNotAllowed":      true,
	"StatusRequestTimeout":        true,
	"StatusConflict":              true,
	"StatusGone":                  true,
	"StatusRequestEntityTooLarge": true,
	"StatusUnsupportedMediaType":  true,
	"StatusUnprocessableEntity":   true,
	"StatusTooManyRequests":       true,
	"StatusInternalServerError":   true,
	"StatusNotImplemented":        true,
	"StatusBadGateway":            true,
	"StatusServiceUnavailable":    true,
	"StatusGatewayTimeout":        true,
}

func ResolveHttpExpr(raw string) (string, error) {
	if _, err := strconv.Atoi(raw); err == nil {
		return raw, nil
	}
	if ValidHttpStatusNames[raw] {
		return "http." + raw, nil
	}
	return "", fmt.Errorf("unknown http status: %q (use a number or net/http constant name like StatusBadRequest)", raw)
}
