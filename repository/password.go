package repository

import "context"

type ForgotPasswordRepository interface {
	CreateForgotPasswordToken(ctx context.Context, email string, token string) error
	GetForgotPasswordToken(ctx context.Context, email string) (string, error)
	DeleteForgotPasswordToken(ctx context.Context, email string) error
}
