package jwt

import (
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
