package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/security"
)

type AccessToken struct {
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	ExpiredAt time.Time `json:"expired_at"`
	JTI       string    `json:"-"`
	UserID    string    `json:"-"`
	SessionID string    `json:"-"`
}

type AccessTokenClaims struct {
	*jwt.StandardClaims
	UserID    string
	SessionID string
}

func CreateAccessToken(userID string, sessionID string) (*AccessToken, error) {
	expiresAt := time.Now().Add(5 * time.Minute)
	jti := string(security.GenerateRandomID())

	log.Info().Msg("creating access token claims")
	claims := AccessTokenClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Id:        jti,
		},
		UserID:    userID,
		SessionID: sessionID,
	}

	log.Info().Msg("creating access token")
	tokenString, err := generateJWTToken(claims)
	if err != nil {
		log.Error().Err(err).Msg("error creating access token")
		return nil, err
	}

	at := &AccessToken{
		Value:     tokenString,
		Type:      "Bearer",
		ExpiredAt: expiresAt,
		JTI:       jti,
		UserID:    userID,
		SessionID: sessionID,
	}
	return at, nil
}

func ParseAccessToken(jwtString string) (*AccessToken, bool, error) {
	claims := &AccessTokenClaims{}

	// parse the token
	t, err := parseJWTToken(jwtString, claims)
	if err != nil {
		log.Error().Err(err).Msg("error parsing access token")
		return nil, false, err
	}

	at := &AccessToken{
		Value:     jwtString,
		Type:      "Bearer",
		ExpiredAt: time.Unix(claims.ExpiresAt, 0),
		JTI:       claims.Id,
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
	}

	return at, t.Valid, nil
}
