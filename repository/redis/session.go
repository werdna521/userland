package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/werdna521/userland/repository"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, s *repository.Session) error
	CreateAccessToken(ctx context.Context,
		at *repository.AccessToken,
		expiresIn time.Duration,
	) error
	CreateRefreshToken(ctx context.Context,
		rt *repository.RefreshToken,
		expiresIn time.Duration,
	) error
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
) error {
	key := r.getSessionKey(s.UserID, s.ID)
	now := time.Now()

	s.CreatedAt = now
	s.UpdatedAt = now

	return r.rdb.HSet(ctx, key, r.toSessionFields(s)).Err()
}

// TODO: make session expire after access token expires
func (r *BaseSessionRepository) AddUserSession(
	ctx context.Context,
	s *repository.Session,
) error {
	key := r.getUserSessionsKey(s.UserID)
	return r.rdb.SAdd(ctx, key, s.ID).Err()
}

func (r *BaseSessionRepository) CreateAccessToken(
	ctx context.Context,
	at *repository.AccessToken,
	expiresIn time.Duration,
) error {
	key := r.getAccessTokenKey(at)
	return r.rdb.SetEX(ctx, key, at.ID, expiresIn).Err()
}

func (r *BaseSessionRepository) CreateRefreshToken(
	ctx context.Context,
	rt *repository.RefreshToken,
	expiresIn time.Duration,
) error {
	key := r.getRefreshTokenKey(rt)
	return r.rdb.SetEX(ctx, key, rt.ID, expiresIn).Err()
}
