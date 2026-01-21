package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type TestAuthenticator struct{}

var secret = "secret"

var testClaims = jwt.MapClaims{
	"iss": "test-aud ",
	"aud": "test-aud",
	"sub": int64(42),
	"exp": time.Now().Add(time.Hour).Unix(),
}

func (t *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	return token.SignedString([]byte(secret)) 
}

func (t *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
}