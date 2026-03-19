package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zakiverse/zakiverse-api/util/trace"
)

type Claims struct {
	AccountId uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	jwt.RegisteredClaims
}

type GenerateParam struct {
	AccountId     uuid.UUID
	Username      string
	Role          string
	Secret        string
	ExpiryMinutes int
}

func Generate(param GenerateParam) (string, error) {
	now := time.Now()
	claims := Claims{
		AccountId: param.AccountId,
		Username:  param.Username,
		Role:      param.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(param.ExpiryMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(param.Secret))
	if err != nil {
		return "", trace.Wrap(err)
	}

	return tokenString, nil
}

func Parse(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, trace.Wrap(err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, trace.Wrap(jwt.ErrSignatureInvalid)
	}

	return claims, nil
}
