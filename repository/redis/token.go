package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/security"
)

type TokenRepository interface {
	CreateForgotPasswordToken(ctx context.Context, userID string, token string) error
	GetForgotPasswordToken(ctx context.Context, token string) (string, error)
	DeleteForgotPasswordToken(ctx context.Context, token string) error
	CreateEmailVerificationToken(ctx context.Context, userID string, token string) error
	GetEmailVerificationToken(ctx context.Context, userID string) (string, error)
	DeleteEmailVerificationToken(ctx context.Context, userID string) error
	CreateEmailChangeToken(ctx context.Context, userID string, t *repository.EmailChangeToken) error
	GetEmailChangeToken(ctx context.Context, userID string) (*repository.EmailChangeToken, error)
	DeleteEmailChangeToken(ctx context.Context, userID string) error
}

type BaseTokenRepository struct {
	rdb *redis.Client
}

func NewBaseTokenRepository(rdb *redis.Client) *BaseTokenRepository {
	return &BaseTokenRepository{
		rdb: rdb,
	}
}

func (r *BaseTokenRepository) getForgotPasswordTokenKey(token string) string {
	return fmt.Sprintf("%s:%s:%s", forgotPasswordKey, tokenKey, token)
}

func (r *BaseTokenRepository) getEmailVerificationTokenKey(userID string) string {
	return fmt.Sprintf("%s:%s:%s:%s", userKey, userID, verificationKey, tokenKey)
}

func (r *BaseTokenRepository) getEmailChangeTokenKey(userID string) string {
	return fmt.Sprintf("%s:%s:%s:%s", userKey, userID, emailChangeVerificationKey, tokenKey)
}

func (r *BaseTokenRepository) CreateForgotPasswordToken(
	ctx context.Context,
	userID string,
	token string,
) error {
	key := r.getForgotPasswordTokenKey(token)
	return r.rdb.SetEX(ctx, key, userID, security.TokenLife).Err()
}

func (r *BaseTokenRepository) GetForgotPasswordToken(
	ctx context.Context,
	token string,
) (string, error) {
	key := r.getForgotPasswordTokenKey(token)

	email, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", repository.NewNotFoundError()
	}

	return email, err
}

func (r *BaseTokenRepository) DeleteForgotPasswordToken(
	ctx context.Context,
	token string,
) error {
	key := r.getForgotPasswordTokenKey(token)

	err := r.rdb.Unlink(ctx, key).Err()
	if err == redis.Nil {
		return repository.NewNotFoundError()
	}

	return err
}

func (r *BaseTokenRepository) CreateEmailVerificationToken(
	ctx context.Context,
	userID string,
	token string,
) error {
	key := r.getEmailVerificationTokenKey(userID)
	return r.rdb.SetEX(ctx, key, token, security.TokenLife).Err()
}

func (r *BaseTokenRepository) GetEmailVerificationToken(
	ctx context.Context,
	userID string,
) (string, error) {
	key := r.getEmailVerificationTokenKey(userID)

	token, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", repository.NewNotFoundError()
	}

	return token, err
}

func (r *BaseTokenRepository) DeleteEmailVerificationToken(
	ctx context.Context,
	userID string,
) error {
	key := r.getEmailVerificationTokenKey(userID)

	err := r.rdb.Unlink(ctx, key).Err()
	if err == redis.Nil {
		return repository.NewNotFoundError()
	}

	return err
}

func (r *BaseTokenRepository) CreateEmailChangeToken(
	ctx context.Context,
	userID string,
	t *repository.EmailChangeToken,
) error {
	key := r.getEmailChangeTokenKey(userID)

	err := r.rdb.HSet(ctx, key, hEmailChangeNewEmailKey, t.NewEmail, hEmailChangeToken, t.Token).Err()
	if err != nil {
		log.Error().Err(err).Msg("failed to store email change request token")
		return err
	}

	err = r.rdb.Expire(ctx, key, security.TokenLife).Err()
	return err
}

func (r *BaseTokenRepository) GetEmailChangeToken(
	ctx context.Context,
	userID string,
) (*repository.EmailChangeToken, error) {
	key := r.getEmailChangeTokenKey(userID)

	res, err := r.rdb.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return nil, repository.NewNotFoundError()
	}

	t := &repository.EmailChangeToken{
		NewEmail: res[hEmailChangeNewEmailKey],
		Token:    res[hEmailChangeToken],
	}

	return t, err
}

func (r *BaseTokenRepository) DeleteEmailChangeToken(
	ctx context.Context,
	userID string,
) error {
	key := r.getEmailChangeTokenKey(userID)

	err := r.rdb.Unlink(ctx, key).Err()
	if err == redis.Nil {
		return repository.NewNotFoundError()
	}

	return err
}
