package trace

import (
	"github.com/pkg/errors"
	runtime "github.com/zakiverse/zakiverse-api/util/runtime"
)

func WrapFunc(message string) string {
	return message + runtime.GetFuncName()
}

func Wrap(err error) error {
	return errors.Wrap(err, runtime.GetFuncName())
}
