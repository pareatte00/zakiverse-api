package trace

import (
	"github.com/pkg/errors"
	"github.com/zakiverse/zakiverse-api/core/cst"
	runtime "github.com/zakiverse/zakiverse-api/util/runtime"
)

func WrapFunc(message string) string {
	return message + runtime.GetFuncName()
}

func Wrap(err error, trace cst.ExceptionTrace) error {
	return errors.Wrap(err, string(trace)+runtime.GetFuncName())
}
