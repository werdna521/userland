package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/repository"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, s *repository.Session, expiresIn time.Duration) error
	AddUserSession(ctx context.Context, s *repository.Session) error
	UpdateSessionExpiryTime(
		ctx context.Context,
		s *repository.Session,
		expiresIn time.Duration,
	) error
	CreateAccessToken(ctx context.Context,
		at *repository.AccessToken,
		expiresIn time.Duration,
	) error
	CheckAccessToken(ctx context.Context, at *repository.AccessToken) (bool, error)
	CreateRefreshToken(ctx context.Context,
		rt *repository.RefreshToken,
		expiresIn time.Duration,
	) error
	CheckRefreshToken(ctx context.Context, rt *repository.RefreshToken) (bool, error)
}

type BaseSessionRepository struct {
	rdb *redis.Client
}

func NewBaseSessionRepository(rdb *redis.Client) *BaseSessionRepository {
	return &BaseSessionRepository{
		rdb: rdb,
	}
}

func (r *BaseSessionRepository) getSessionKey(userID string, sessionID string) string {
	return fmt.Sprintf("%s:%s:%s:%s", userKey, userID, sessionKey, sessionID)
}

func (r *BaseSessionRepository) getUserSessionsKey(userID string) string {
	return fmt.Sprintf("%s:%s:%s", userKey, userID, sessionKey)
}

func (r *BaseSessionRepository) getAccessTokenKey(at *repository.AccessToken) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", userKey, at.UserID, sessionKey, at.SessionID, accessTokenKey)
}

func (r *BaseSessionRepository) getRefreshTokenKey(rt *repository.RefreshToken) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", userKey, rt.UserID, sessionKey, rt.SessionID, refreshTokenKey)
}

func (r *BaseSessionRepository) toSessionFields(s *repository.Session) map[string]interface{} {
	return map[string]interface{}{
		hSessionIDKey:        s.ID,
		hSessionClientKey:    s.Client,
		hSessionCreatedAtKey: s.CreatedAt,
		hSessionUpdatedAtKey: s.UpdatedAt,
	}
}

// TODO: make session expire after access token expires
func (r *BaseSessionRepository) CreateSession(
	ctx context.Context,
	s *repository.Session,
	expiresIn time.Duration,
) error {
	key := r.getSessionKey(s.UserID, s.ID)
	now := time.Now()

	s.CreatedAt = now
	s.UpdatedAt = now

	err := r.rdb.HSet(ctx, key, r.toSessionFields(s)).Err()
	if err != nil {
		log.Error().Err(err).Msg("failed to create session")
		return err
	}

	err = r.rdb.Expire(ctx, key, expiresIn).Err()
	return err
}

// TODO: make session expire after access token expires
func (r *BaseSessionRepository) AddUserSession(
	ctx context.Context,
	s *repository.Session,
) error {
	key := r.getUserSessionsKey(s.UserID)
	return r.rdb.SAdd(ctx, key, s.ID).Err()
}

func (r *BaseSessionRepository) UpdateSessionExpiryTime(
	ctx context.Context,
	s *repository.Session,
	expiresIn time.Duration,
) error {
	key := r.getSessionKey(s.UserID, s.ID)
	now := time.Now()

	exp, err := r.rdb.TTL(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to get session expiry time")
		return err
	}

	err = r.rdb.HSet(ctx, key, hSessionUpdatedAtKey, now).Err()
	if err != nil {
		log.Error().Err(err).Msg("failed to touch session")
		return err
	}

	// we only update the session expiry time if it is less than the new expiry time
	if exp.Seconds() < expiresIn.Seconds() {
		err = r.rdb.Expire(ctx, key, expiresIn).Err()
	}
	return err
}

func (r *BaseSessionRepository) CreateAccessToken(
	ctx context.Context,
	at *repository.AccessToken,
	expiresIn time.Duration,
) error {
	key := r.getAccessTokenKey(at)
	return r.rdb.SetEX(ctx, key, at.ID, expiresIn).Err()
}

func (r *BaseSessionRepository) CheckAccessToken(
	ctx context.Context,
	at *repository.AccessToken,
) (bool, error) {
	key := r.getAccessTokenKey(at)

	jti, err := r.rdb.Get(ctx, key).Result()
	if jti == "" {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return jti == at.ID, nil
}

func (r *BaseSessionRepository) CreateRefreshToken(
	ctx context.Context,
	rt *repository.RefreshToken,
	expiresIn time.Duration,
) error {
	key := r.getRefreshTokenKey(rt)
	return r.rdb.SetEX(ctx, key, rt.ID, expiresIn).Err()
}

func (r *BaseSessionRepository) CheckRefreshToken(
	ctx context.Context,
	rt *repository.RefreshToken,
) (bool, error) {
	key := r.getRefreshTokenKey(rt)

	jti, err := r.rdb.Get(ctx, key).Result()
	if jti == "" {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return jti == rt.ID, nil
}
