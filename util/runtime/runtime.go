package runtime

import (
	"fmt"
	"runtime"
)

const (
	formatFunctionName  = "[%s]"
	unknownFunctionName = "unknown"
)

func GetFuncName() string {
	pc, _, _, success := runtime.Caller(2)
	if !success {
		return fmt.Sprintf(formatFunctionName, unknownFunctionName)
	}
	return fmt.Sprintf(formatFunctionName, runtime.FuncForPC(pc).Name())
}
