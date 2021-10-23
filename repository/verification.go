package repository

import "context"

type VerificationCode string

type EmailVerificationRepository interface {
	CreateVerification(ctx context.Context, email string, vc VerificationCode) error
}
