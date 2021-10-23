package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/werdna521/userland/repository"
)

const (
	userKey         = "users"
	verificationKey = "verification"
)

type EmailVerificationRepository struct {
	rdb *redis.Client
}

func NewVerificationRepository(rdb *redis.Client) *EmailVerificationRepository {
	return &EmailVerificationRepository{
		rdb: rdb,
	}
}

func (vr *EmailVerificationRepository) getKey(email string) string {
	return fmt.Sprintf("%s:%s:%s", userKey, email, verificationKey)
}

func (vr *EmailVerificationRepository) CreateVerification(
	ctx context.Context,
	email string,
	verificationCode repository.VerificationCode,
) error {
	key := vr.getKey(email)
	return vr.rdb.SetEX(ctx, key, string(verificationCode), 60*time.Second).Err()
}
