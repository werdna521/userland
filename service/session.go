package service

import (
	"context"

	"github.com/rs/zerolog/log"
	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/redis"
	"github.com/werdna521/userland/security/jwt"
	"github.com/werdna521/userland/utils/slice"
)

type SessionService interface {
	GenerateRefreshToken(ctx context.Context, at *jwt.AccessToken) (*jwt.RefreshToken, e.Error)
	GenerateAccessToken(ctx context.Context, rt *jwt.RefreshToken) (*jwt.AccessToken, e.Error)
	ListSessions(ctx context.Context, at *jwt.AccessToken) ([]*repository.Session, e.Error)
	RemoveSession(ctx context.Context, session *repository.Session) e.Error
	RemoveAllOtherSessions(ctx context.Context, session *repository.Session) e.Error
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

func (s *BaseSessionService) ListSessions(
	ctx context.Context,
	at *jwt.AccessToken,
) ([]*repository.Session, e.Error) {
	log.Info().Msg("getting all active sessions")
	sessions, err := s.sr.GetAllSessions(ctx, at.UserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all active sessions")
		return nil, e.NewInternalServerError()
	}

	return sessions, nil
}

func (s *BaseSessionService) RemoveSession(
	ctx context.Context,
	session *repository.Session,
) e.Error {
	log.Info().Msg("removing session from redis")
	err := s.sr.DeleteSession(ctx, session)
	if err != nil {
		log.Error().Err(err).Msg("failed to remove session from redis")
		return e.NewInternalServerError()
	}

	accessToken := &repository.AccessToken{
		UserID:    session.UserID,
		SessionID: session.ID,
	}
	log.Info().Msg("revoking access token")
	err = s.sr.DeleteAccessToken(ctx, accessToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to revoke access token")
		return e.NewInternalServerError()
	}

	refreshToken := &repository.RefreshToken{
		UserID:    session.UserID,
		SessionID: session.ID,
	}
	log.Info().Msg("revoking refresh token")
	err = s.sr.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to revoke refresh token")
		return e.NewInternalServerError()
	}

	log.Info().Msg("removing session id from index")
	err = s.sr.RemoveUserSessionFromIndex(ctx, session.UserID, session.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to remove session id from index")
		return e.NewInternalServerError()
	}

	return nil
}

func (s *BaseSessionService) RemoveAllOtherSessions(
	ctx context.Context,
	session *repository.Session,
) e.Error {
	log.Info().Msg("getting all active sessions")
	sessions, err := s.sr.GetAllSessions(ctx, session.UserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all active sessions")
		return e.NewInternalServerError()
	}

	log.Info().Msg("removing current session from sessions list")
	sessions = slice.FilterSession(sessions, func(s *repository.Session) bool {
		return s.ID != session.ID
	})

	log.Info().Msg("removing all other sessions from redis")
	for _, session := range sessions {
		err = s.RemoveSession(ctx, session)
		if err != nil {
			log.Error().Err(err).Msgf("failed to remove session %s", session.ID)
			return e.NewInternalServerError()
		}
	}

	return nil
}
