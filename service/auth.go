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
	VerifyEmail(ctx context.Context, email string, code string) e.Error
}

type BaseAuthService struct {
	ur  repository.UserRepository
	evr repository.EmailVerificationRepository
}

func NewBaseAuthService(
	ur repository.UserRepository,
	evr repository.EmailVerificationRepository,
) *BaseAuthService {
	return &BaseAuthService{
		ur:  ur,
		evr: evr,
	}
}

func (bas *BaseAuthService) Register(ctx context.Context, u *repository.User) e.Error {
	log.Info().Msg("checking if email is available")
	_, err := bas.ur.GetUserByEmail(ctx, u.Email)
	if err == nil {
		log.Error().Err(err).Msg("email is not available")
		return e.NewConflictError("email is not available")
	}
	if _, ok := err.(repository.NotFoundError); !ok {
		log.Error().Err(err).Msg("fail to check email availability")
		return e.NewInternalServerError()
	}

	log.Info().Msg("hashing password")
	hash, err := security.HashPassword(u.Password)
	if err != nil {
		log.Error().Err(err).Msg("fail to hash password")
		return e.NewInternalServerError()
	}
	u.Password = hash

	log.Info().Msg("creating and registering user")
	err = bas.ur.CreateUser(ctx, u)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return e.NewInternalServerError()
	}

	log.Info().Msg("generating verification code")
	verificationCode := security.GenerateRandomID()

	log.Info().Msg("storing email verification code")
	err = bas.evr.CreateVerification(ctx, u.Email, string(verificationCode))
	if err != nil {
		log.Error().Err(err).Msg("failed to store email verification code")
		return e.NewInternalServerError()
	}

	// TODO: send email with verification code/link
	log.Debug().Msgf("verification code: %s", verificationCode)

	return nil
}

func (bas *BaseAuthService) VerifyEmail(
	ctx context.Context,
	email string,
	verificationCode string,
) e.Error {
	// TODO: to make it even safer, also check for the user's existence in the
	// database. ideally, this case will be impossible, though.

	log.Info().Msg("retrieving verification code from repository")
	storedCode, err := bas.evr.GetVerification(ctx, email)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("verification details not found")
		return e.NewNotFoundError("account not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get verification code")
		return e.NewInternalServerError()
	}

	log.Info().Msg("checking verification code")
	if storedCode != verificationCode {
		log.Error().Msg("invalid verification code")
		return e.NewUnauthorizedError("invalid verification code")
	}

	log.Info().Msg("activating user account")
	err = bas.ur.UpdateUserActivationStatusByEmail(ctx, email, true)
	if err != nil {
		log.Error().Err(err).Msg("failed to activate user account")
		return e.NewInternalServerError()
	}

	// still, remove the data from redis, even though it doesn't really matter tbh
	log.Info().Msg("removing verification details from redis")
	err = bas.evr.DeleteVerification(ctx, email)
	if _, ok := err.(repository.NotFoundError); !ok {
		log.Error().Err(err).Msg("failed to remove verification details from redis")
		return e.NewInternalServerError()
	}

	return nil
}
