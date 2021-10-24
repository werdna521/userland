package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/werdna521/userland/repository"
)

const (
	userKey           = "user"
	verificationKey   = "verification"
	forgotPasswordKey = "forgotPassword"
	tokenKey          = "token"
)

type TokenRepository interface {
	CreateForgotPasswordToken(ctx context.Context, email string, token string) error
	GetForgotPasswordToken(ctx context.Context, token string) (string, error)
	DeleteForgotPasswordToken(ctx context.Context, token string) error
	CreateEmailVerificationToken(ctx context.Context, email string, token string) error
	GetEmailVerificationToken(ctx context.Context, email string) (string, error)
	DeleteEmailVerificationToken(ctx context.Context, email string) error
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

func (r *BaseTokenRepository) getEmailVerificationTokenKey(email string) string {
	return fmt.Sprintf("%s:%s:%s:%s", userKey, email, verificationKey, tokenKey)
}

func (r *BaseTokenRepository) CreateForgotPasswordToken(
	ctx context.Context,
	email string,
	token string,
) error {
	key := r.getForgotPasswordTokenKey(token)
	return r.rdb.SetEX(ctx, key, email, 5*time.Minute).Err()
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
	email string,
) error {
	key := r.getForgotPasswordTokenKey(email)

	err := r.rdb.Del(ctx, key).Err()
	if err == redis.Nil {
		return repository.NewNotFoundError()
	}

	return err
}

func (r *BaseTokenRepository) CreateEmailVerificationToken(
	ctx context.Context,
	email string,
	token string,
) error {
	key := r.getEmailVerificationTokenKey(email)
	return r.rdb.SetEX(ctx, key, token, 5*time.Minute).Err()
}

func (r *BaseTokenRepository) GetEmailVerificationToken(
	ctx context.Context,
	email string,
) (string, error) {
	key := r.getEmailVerificationTokenKey(email)

	token, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", repository.NewNotFoundError()
	}

	return token, err
}

func (r *BaseTokenRepository) DeleteEmailVerificationToken(
	ctx context.Context,
	email string,
) error {
	key := r.getEmailVerificationTokenKey(email)

	err := r.rdb.Del(ctx, key).Err()
	if err == redis.Nil {
		return repository.NewNotFoundError()
	}

	return err
}