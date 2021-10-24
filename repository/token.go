package repository

import "context"

type TokenRepository interface {
	CreateForgotPasswordToken(ctx context.Context, email string, token string) error
	GetForgotPasswordToken(ctx context.Context, token string) (string, error)
	DeleteForgotPasswordToken(ctx context.Context, token string) error
	CreateEmailVerificationToken(ctx context.Context, email string, token string) error
	GetEmailVerificationToken(ctx context.Context, email string) (string, error)
	DeleteEmailVerificationToken(ctx context.Context, email string) error
}
