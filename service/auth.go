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
	SendEmailVerification(ctx context.Context, email string) e.Error
	VerifyEmail(ctx context.Context, email string, code string) e.Error
	ForgotPassword(ctx context.Context, email string) e.Error
}

type BaseAuthService struct {
	ur  repository.UserRepository
	evr repository.EmailVerificationRepository
	pvr repository.ForgotPasswordRepository
}

func NewBaseAuthService(
	ur repository.UserRepository,
	evr repository.EmailVerificationRepository,
	pvr repository.ForgotPasswordRepository,
) *BaseAuthService {
	return &BaseAuthService{
		ur:  ur,
		evr: evr,
		pvr: pvr,
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

	log.Info().Msg("starting send verification email service")
	err = bas.SendEmailVerification(ctx, u.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to send verification email")
		return err.(e.Error)
	}

	return nil
}

func (bas *BaseAuthService) SendEmailVerification(
	ctx context.Context,
	email string,
) e.Error {
	log.Info().Msg("retrieving user from database")
	u, err := bas.ur.GetUserByEmail(ctx, email)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("user not found")
		return e.NewNotFoundError("user not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		return e.NewInternalServerError()
	}

	if u.IsActive {
		log.Error().Msg("user is already active")
		return e.NewBadRequestError("user is already active")
	}

	log.Info().Msg("generating verification token")
	verificationToken := security.GenerateRandomID()

	log.Info().Msg("storing email verification token")
	err = bas.evr.CreateVerificationToken(ctx, email, string(verificationToken))
	if err != nil {
		log.Error().Err(err).Msg("failed to store email verification token")
		return e.NewInternalServerError()
	}

	// TODO: send email with verification token/link
	log.Debug().Msgf("verification token: %s", verificationToken)

	return nil
}

func (bas *BaseAuthService) VerifyEmail(
	ctx context.Context,
	email string,
	verificationToken string,
) e.Error {
	// TODO: to make it even safer, also check for the user's existence in the
	// database. ideally, this case will be impossible, though.

	log.Info().Msg("retrieving verification token from repository")
	storedToken, err := bas.evr.GetVerificationToken(ctx, email)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("verification token not found")
		return e.NewNotFoundError("user not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get verification token")
		return e.NewInternalServerError()
	}

	log.Info().Msg("checking verification verification token")
	if storedToken != verificationToken {
		log.Error().Msg("invalid verification token")
		return e.NewUnauthorizedError("invalid verification token")
	}

	log.Info().Msg("activating user account")
	err = bas.ur.UpdateUserActivationStatusByEmail(ctx, email, true)
	if err != nil {
		log.Error().Err(err).Msg("failed to activate user account")
		return e.NewInternalServerError()
	}

	// still, remove the data from redis, even though it doesn't really matter tbh
	log.Info().Msg("removing verification details from redis")
	err = bas.evr.DeleteVerificationToken(ctx, email)
	if _, ok := err.(repository.NotFoundError); !ok && err != nil {
		log.Error().Err(err).Msg("failed to remove verification details from redis")
		return e.NewInternalServerError()
	}

	return nil
}

func (bas *BaseAuthService) ForgotPassword(ctx context.Context, email string) e.Error {
	log.Info().Msg("retrieving user from the db")
	u, err := bas.ur.GetUserByEmail(ctx, email)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("user not found")
		return e.NewNotFoundError("user not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve user")
		return e.NewInternalServerError()
	}

	log.Info().Msg("checking user activation status")
	if !u.IsActive {
		log.Error().Msg("user is not active")
		return e.NewBadRequestError("user is not active")
	}

	log.Info().Msg("generating forgot password token")
	token := security.GenerateRandomID()

	log.Info().Msg("storing forgot password token")
	err = bas.pvr.CreateForgotPasswordToken(ctx, email, string(token))
	if err != nil {
		log.Error().Err(err).Msg("failed to store forgot password token")
		return e.NewInternalServerError()
	}

	// TODO: send email with verification token/link
	log.Debug().Msgf("forgot password token: %s", token)

	return nil
}
