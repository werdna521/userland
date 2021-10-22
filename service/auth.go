package service

import (
	"context"

	"github.com/rs/zerolog/log"
	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/security"
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
	log.Info().Msg("checking if email is available")
	_, err := bas.ur.GetUserByEmail(ctx, u.Email)
	if err == nil {
		log.Error().Stack().Err(err).Msg("email is not available")
		return e.NewConflictError("email is not available")
	}
	if _, ok := err.(repository.NotFoundError); !ok {
		log.Error().Stack().Err(err).Msg("fail to check email availability")
		return e.NewInternalServerError()
	}

	// TODO: hash password
	log.Info().Msg("hashing password")
	hash, err := security.HashPassword(u.Password)
	if err != nil {
		log.Error().Stack().Err(err).Msg("fail to hash password")
		return e.NewInternalServerError()
	}
	u.Password = hash

	log.Info().Msg("creating and registering user")
	err = bas.ur.CreateUser(ctx, u)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to create user")
		return e.NewInternalServerError()
	}

	// TODO: send email via SendVerification service handler

	return nil
}
