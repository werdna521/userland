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
	GetSession(ctx context.Context, userID string, sessionID string) (*repository.Session, error)
	GetAllSessions(ctx context.Context, userID string) ([]*repository.Session, error)
	DeleteSession(ctx context.Context, s *repository.Session) error
	AddUserSessionToIndex(ctx context.Context, s *repository.Session) error
	RemoveUserSessionFromIndex(ctx context.Context, userID string, sessionID string) error
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
	DeleteAccessToken(ctx context.Context, at *repository.AccessToken) error
	CreateRefreshToken(ctx context.Context,
		rt *repository.RefreshToken,
		expiresIn time.Duration,
	) error
	CheckRefreshToken(ctx context.Context, rt *repository.RefreshToken) (bool, error)
	DeleteRefreshToken(ctx context.Context, rt *repository.RefreshToken) error
}

type BaseSessionRepository struct {
	rdb *redis.Client
}

func NewBaseSessionRepository(rdb *redis.Client) *BaseSessionRepository {
	return &BaseSessionRepository{
		rdb: rdb,
	}
}

// TODO: get keys shouldn't accept struct as params
func (r *BaseSessionRepository) getSessionKey(userID string, sessionID string) string {
	return fmt.Sprintf("%s:%s:%s:%s", userKey, userID, sessionKey, sessionID)
}

func (r *BaseSessionRepository) getSessionIndexKey(userID string) string {
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
		hSessionClientKey:    s.Client,
		hSessionIPAddress:    s.IPAddress,
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

func (r *BaseSessionRepository) GetSession(
	ctx context.Context,
	userID string,
	sessionID string,
) (*repository.Session, error) {
	key := r.getSessionKey(userID, sessionID)

	res, err := r.rdb.HGetAll(ctx, key).Result()
	if len(res) == 0 {
		return nil, repository.NewNotFoundError()
	}
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, res[hSessionCreatedAtKey])
	if err != nil {
		log.Error().Err(err).Msg("failed to parse created_at timestamp")
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339, res[hSessionUpdatedAtKey])
	if err != nil {
		log.Error().Err(err).Msg("failed to parse updated_at timestamp")
		return nil, err
	}

	session := &repository.Session{
		ID:        sessionID,
		UserID:    userID,
		Client:    res[hSessionClientKey],
		IPAddress: res[hSessionIPAddress],
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return session, nil
}

func (r *BaseSessionRepository) GetAllSessions(
	ctx context.Context,
	userID string,
) ([]*repository.Session, error) {
	sessionIndexKey := r.getSessionIndexKey(userID)

	sessionIDs, err := r.rdb.SMembers(ctx, sessionIndexKey).Result()
	if err != nil {
		return nil, err
	}

	sessions := []*repository.Session{}
	for _, sessionID := range sessionIDs {
		session, err := r.GetSession(ctx, userID, sessionID)

		// if session is not found, this means the session is expired, hence we
		// remove it from the index
		if _, ok := err.(repository.NotFoundError); ok {
			err = r.RemoveUserSessionFromIndex(ctx, userID, sessionID)
			if err != nil {
				log.Error().Err(err).Msg("failed to remove session from index")
				return nil, err
			}
			continue
		}
		if err != nil {
			return nil, err
		}

		// else, we append to the slice
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *BaseSessionRepository) DeleteSession(
	ctx context.Context,
	s *repository.Session,
) error {
	key := r.getSessionKey(s.UserID, s.ID)
	return r.rdb.Unlink(ctx, key).Err()
}

// a tiny problem with redis: we can't set expiration time for a single element
// in a set. we'll have to handle deletion manually in the code :(
func (r *BaseSessionRepository) AddUserSessionToIndex(
	ctx context.Context,
	s *repository.Session,
) error {
	key := r.getSessionIndexKey(s.UserID)
	return r.rdb.SAdd(ctx, key, s.ID).Err()
}

func (r *BaseSessionRepository) RemoveUserSessionFromIndex(
	ctx context.Context,
	userID string,
	sessionID string,
) error {
	key := r.getSessionIndexKey(userID)
	return r.rdb.SRem(ctx, key, sessionID).Err()
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

func (r *BaseSessionRepository) DeleteAccessToken(
	ctx context.Context,
	at *repository.AccessToken,
) error {
	key := r.getAccessTokenKey(at)
	return r.rdb.Unlink(ctx, key).Err()
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

func (r *BaseSessionRepository) DeleteRefreshToken(
	ctx context.Context,
	rt *repository.RefreshToken,
) error {
	key := r.getRefreshTokenKey(rt)
	return r.rdb.Unlink(ctx, key).Err()
}
