package jwt

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

type jwtConfig struct {
	secret []byte
}

var cfg jwtConfig

func init() {
	cfg = jwtConfig{
		secret: []byte(os.Getenv("JWT_SECRET")),
	}
}

func generateJWTToken(claims jwt.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(cfg.secret)
}

func parseJWTToken(jwtString string, claims jwt.Claims) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(jwtString, claims, func(t *jwt.Token) (interface{}, error) {
		// check token signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}

		return cfg.secret, nil
	})
	if _, ok := err.(*jwt.ValidationError); ok {
		return nil, NewInvalidTokenError()
	}
	return t, err
}
