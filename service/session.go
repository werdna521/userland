package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/redis"
	"github.com/werdna521/userland/security/jwt"
)

type SessionService interface {
	GenerateRefreshToken(ctx context.Context, at *jwt.AccessToken) (*jwt.RefreshToken, e.Error)
}

type BaseSessionService struct {
	sr redis.SessionRepository
}

func NewBaseSessionService(sr redis.SessionRepository) *BaseSessionService {
	return &BaseSessionService{
		sr: sr,
	}
}

func (s *BaseSessionService) GenerateRefreshToken(
	ctx context.Context,
	at *jwt.AccessToken,
) (*jwt.RefreshToken, e.Error) {
	log.Info().Msg("generating refresh token")
	rt, err := jwt.CreateRefreshToken(at.UserID, at.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate refresh token")
	}

	log.Info().Msg("storing refresh token in redis")
	token := &repository.RefreshToken{
		ID:        rt.JTI,
		UserID:    rt.UserID,
		SessionID: rt.SessionID,
	}
	// TODO: make expiration longer, current duration is just for debugging purposes
	err = s.sr.CreateRefreshToken(ctx, token, 10*time.Minute)
	if err != nil {
		log.Error().Err(err).Msg("failed to store refresh token in redis")
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("updating session expiry time")
	session := &repository.Session{
		ID:     at.SessionID,
		UserID: at.UserID,
	}
	err = s.sr.UpdateSessionExpiryTime(ctx, session, 10*time.Minute)
	if err != nil {
		log.Error().Err(err).Msg("failed to update session expiry time")
		return nil, e.NewInternalServerError()
	}

	return rt, nil
}
