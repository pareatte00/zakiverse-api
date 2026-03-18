package password

import (
	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/util/trace"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string, cost int) (hashPasswordText string, err error) {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", trace.Wrap(err, cst.ExceptionTraceUtils)
	}

	return string(byteHash), nil
}
