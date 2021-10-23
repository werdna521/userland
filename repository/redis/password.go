package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/werdna521/userland/repository"
)

type ForgotPasswordRepository struct {
	rdb *redis.Client
}

func NewForgotPasswordRepository(rdb *redis.Client) *ForgotPasswordRepository {
	return &ForgotPasswordRepository{
		rdb: rdb,
	}
}

func (fpr *ForgotPasswordRepository) getKey(email string) string {
	return fmt.Sprintf("%s:%s:%s:%s", userKey, email, forgotPasswordKey, tokenKey)
}

func (fpr *ForgotPasswordRepository) CreateForgotPasswordToken(
	ctx context.Context,
	email string,
	token string,
) error {
	key := fpr.getKey(email)
	return fpr.rdb.SetEX(ctx, key, token, 5*time.Minute).Err()
}

func (fpr *ForgotPasswordRepository) GetForgotPasswordToken(
	ctx context.Context,
	email string,
) (string, error) {
	key := fpr.getKey(email)

	token, err := fpr.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", repository.NewNotFoundError()
	}

	return token, err
}

func (fpr *ForgotPasswordRepository) DeleteForgotPasswordToken(
	ctx context.Context,
	email string,
) error {
	key := fpr.getKey(email)

	err := fpr.rdb.Del(ctx, key).Err()
	if err == redis.Nil {
		return repository.NewNotFoundError()
	}

	return err
}
