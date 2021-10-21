package service

import (
	"context"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/repository"
)

type AuthService interface {
	Register(ctx context.Context, user *repository.User) e.Error
}

type BaseAuthService struct {
	ur repository.UserRepository
}

func NewBaseAuthService(ur repository.UserRepository) *BaseAuthService {
	return &BaseAuthService{
		ur: ur,
	}
}

func (bas *BaseAuthService) Register(ctx context.Context, u *repository.User) e.Error {
	// TODO: check if user can be crated, ie. check if email is unique
	// TODO: hash password

	err := bas.ur.CreateUser(ctx, u)
	if err != nil {
		return e.NewInternalServerError()
	}

	// TODO: send email via SendVerification service handler

	return nil
}
