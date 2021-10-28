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
	expiresAt := time.Now().Add(RefreshTokenLife)
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

func ParseRefreshToken(jwtString string) (*RefreshToken, bool, error) {
	claims := &RefreshTokenClaims{}

	// parse the token
	t, err := parseJWTToken(jwtString, claims)
	if err != nil {
		log.Error().Err(err).Msg("error parsing access token")
		return nil, false, err
	}

	rt := &RefreshToken{
		Value:     jwtString,
		Type:      "Bearer",
		ExpiredAt: time.Unix(claims.ExpiresAt, 0),
		JTI:       claims.Id,
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
	}

	return rt, t.Valid, nil
}
