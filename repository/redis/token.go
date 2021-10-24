package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/werdna521/userland/repository"
)

type TokenRepository struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) *TokenRepository {
	return &TokenRepository{
		rdb: rdb,
	}
}

func (fpr *TokenRepository) getForgotPasswordTokenKey(token string) string {
	return fmt.Sprintf("%s:%s:%s", forgotPasswordKey, tokenKey, token)
}

func (evr *TokenRepository) getEmailVerificationTokenKey(email string) string {
	return fmt.Sprintf("%s:%s:%s:%s", userKey, email, verificationKey, tokenKey)
}

func (fpr *TokenRepository) CreateForgotPasswordToken(
	ctx context.Context,
	email string,
	token string,
) error {
	key := fpr.getForgotPasswordTokenKey(token)
	return fpr.rdb.SetEX(ctx, key, email, 5*time.Minute).Err()
}

func (fpr *TokenRepository) GetForgotPasswordToken(
	ctx context.Context,
	token string,
) (string, error) {
	key := fpr.getForgotPasswordTokenKey(token)

	email, err := fpr.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", repository.NewNotFoundError()
	}

	return email, err
}

func (fpr *TokenRepository) DeleteForgotPasswordToken(
	ctx context.Context,
	email string,
) error {
	key := fpr.getForgotPasswordTokenKey(email)

	err := fpr.rdb.Del(ctx, key).Err()
	if err == redis.Nil {
		return repository.NewNotFoundError()
	}

	return err
}

func (evr *TokenRepository) CreateEmailVerificationToken(
	ctx context.Context,
	email string,
	token string,
) error {
	key := evr.getEmailVerificationTokenKey(email)
	return evr.rdb.SetEX(ctx, key, token, 5*time.Minute).Err()
}

func (evr *TokenRepository) GetEmailVerificationToken(
	ctx context.Context,
	email string,
) (string, error) {
	key := evr.getEmailVerificationTokenKey(email)

	token, err := evr.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", repository.NewNotFoundError()
	}

	return token, err
}

func (evr *TokenRepository) DeleteEmailVerificationToken(
	ctx context.Context,
	email string,
) error {
	key := evr.getEmailVerificationTokenKey(email)

	err := evr.rdb.Del(ctx, key).Err()
	if err == redis.Nil {
		return repository.NewNotFoundError()
	}

	return err
}
