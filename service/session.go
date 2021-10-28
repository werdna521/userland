package service

import (
	"context"

	"github.com/rs/zerolog/log"
	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/redis"
	"github.com/werdna521/userland/security/jwt"
)

type SessionService interface {
	GenerateRefreshToken(ctx context.Context, at *jwt.AccessToken) (*jwt.RefreshToken, e.Error)
	GenerateAccessToken(ctx context.Context, rt *jwt.RefreshToken) (*jwt.AccessToken, e.Error)
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
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("storing refresh token in redis")
	token := &repository.RefreshToken{
		ID:        rt.JTI,
		UserID:    rt.UserID,
		SessionID: rt.SessionID,
	}
	err = s.sr.CreateRefreshToken(ctx, token, jwt.RefreshTokenLife)
	if err != nil {
		log.Error().Err(err).Msg("failed to store refresh token in redis")
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("updating session expiry time")
	session := &repository.Session{
		ID:     rt.SessionID,
		UserID: rt.UserID,
	}
	err = s.sr.UpdateSessionExpiryTime(ctx, session, jwt.RefreshTokenLife)
	if err != nil {
		log.Error().Err(err).Msg("failed to update session expiry time")
		return nil, e.NewInternalServerError()
	}

	return rt, nil
}

func (s *BaseSessionService) GenerateAccessToken(
	ctx context.Context,
	rt *jwt.RefreshToken,
) (*jwt.AccessToken, e.Error) {
	log.Info().Msg("generating access token")
	at, err := jwt.CreateAccessToken(rt.UserID, rt.SessionID)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate access token")
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("storing access token in redis")
	token := &repository.AccessToken{
		ID:        at.JTI,
		UserID:    at.UserID,
		SessionID: at.SessionID,
	}
	err = s.sr.CreateAccessToken(ctx, token, jwt.AccessTokenLife)
	if err != nil {
		log.Error().Err(err).Msg("failed to store access token in redis")
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("updating session expiry time")
	session := &repository.Session{
		ID:     at.SessionID,
		UserID: at.UserID,
	}
	err = s.sr.UpdateSessionExpiryTime(ctx, session, jwt.AccessTokenLife)
	if err != nil {
		log.Error().Err(err).Msg("failed to update session expiry time")
		return nil, e.NewInternalServerError()
	}

	return at, nil
}
