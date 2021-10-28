package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/security"
)

type RefreshToken struct {
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	ExpiredAt time.Time `json:"expired_at"`
	JTI       string    `json:"-"`
	UserID    string    `json:"-"`
	SessionID string    `json:"-"`
}

type RefreshTokenClaims struct {
	*jwt.StandardClaims
	UserID    string
	SessionID string
}

func CreateRefreshToken(userID string, sessionID string) (*RefreshToken, error) {
	// TODO: make expiration longer, current duration is just for debugging purposes
	expiresAt := time.Now().Add(10 * time.Minute)
	jti := string(security.GenerateRandomID())

	claims := RefreshTokenClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Id:        jti,
		},
		UserID:    userID,
		SessionID: sessionID,
	}

	log.Info().Msg("creating refresh token")
	tokenString, err := generateJWTToken(claims)
	if err != nil {
		log.Error().Err(err).Msg("failed to create refresh token")
		return nil, err
	}

	rt := &RefreshToken{
		Value:     tokenString,
		Type:      "Bearer",
		ExpiredAt: expiresAt,
		JTI:       jti,
		UserID:    userID,
		SessionID: sessionID,
	}
	return rt, nil
}
