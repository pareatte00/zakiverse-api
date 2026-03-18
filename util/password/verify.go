package password

import "golang.org/x/crypto/bcrypt"

func Verify(passwordText string, passwordHash string) (isMatch bool) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordText))
	return err == nil
}
