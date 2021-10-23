package repository

import "context"

type EmailVerificationRepository interface {
	CreateVerificationToken(ctx context.Context, email string, token string) error
	GetVerificationToken(ctx context.Context, email string) (string, error)
	DeleteVerificationToken(ctx context.Context, email string) error
}
