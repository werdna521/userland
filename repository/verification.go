package repository

import "context"

type EmailVerificationRepository interface {
	CreateVerification(ctx context.Context, email string, vc string) error
	GetVerification(ctx context.Context, email string) (string, error)
	DeleteVerification(ctx context.Context, email string) error
}
