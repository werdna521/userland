package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/werdna521/userland/repository"
)

type EmailVerificationRepository struct {
	rdb *redis.Client
}

func NewVerificationRepository(rdb *redis.Client) *EmailVerificationRepository {
	return &EmailVerificationRepository{
		rdb: rdb,
	}
}

func (evr *EmailVerificationRepository) getKey(email string) string {
	return fmt.Sprintf("%s:%s:%s:%s", userKey, email, verificationKey, tokenKey)
}

func (evr *EmailVerificationRepository) CreateVerificationToken(
	ctx context.Context,
	email string,
	token string,
) error {
	key := evr.getKey(email)
	return evr.rdb.SetEX(ctx, key, token, 5*time.Minute).Err()
}

func (evr *EmailVerificationRepository) GetVerificationToken(
	ctx context.Context,
	email string,
) (string, error) {
	key := evr.getKey(email)

	token, err := evr.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", repository.NewNotFoundError()
	}

	return token, err
}

func (evr *EmailVerificationRepository) DeleteVerificationToken(
	ctx context.Context,
	email string,
) error {
	key := evr.getKey(email)

	err := evr.rdb.Del(ctx, key).Err()
	if err == redis.Nil {
		return repository.NewNotFoundError()
	}

	return err
}
