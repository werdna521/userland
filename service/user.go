package service

import (
	"context"

	"github.com/rs/zerolog/log"
	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/postgres"
	"github.com/werdna521/userland/repository/redis"
	"github.com/werdna521/userland/security"
	"github.com/werdna521/userland/utils/slice"
)

type UserService interface {
	GetInfoDetail(ctx context.Context, userID string) (*repository.UserBio, e.Error)
	UpdateBasicInfo(ctx context.Context, userID string, ub *repository.UserBio) e.Error
	GetCurrentEmail(ctx context.Context, userID string) (string, e.Error)
	RequestEmailChange(ctx context.Context, userID string, newEmail string) e.Error
	VerifyEmailChange(ctx context.Context, userID string, token string) e.Error
	ChangePassword(
		ctx context.Context,
		userID string,
		currentPassword string,
		newPassword string,
	) e.Error
}

type BaseUserService struct {
	ur  postgres.UserRepository
	phr postgres.PasswordHistoryRepository
	tr  redis.TokenRepository
}

func NewBaseUserService(
	ur postgres.UserRepository,
	phr postgres.PasswordHistoryRepository,
	tr redis.TokenRepository,
) *BaseUserService {
	return &BaseUserService{
		ur:  ur,
		phr: phr,
		tr:  tr,
	}
}

func (s *BaseUserService) GetInfoDetail(
	ctx context.Context,
	userID string,
) (*repository.UserBio, e.Error) {
	log.Info().Msg("getting user bio from database")
	ub, err := s.ur.GetUserBioByID(ctx, userID)
	if _, ok := err.(repository.NotFoundError); ok {
		// this shouldn't happen in real-world scenario due to the fact that userID
		// is coming from the access token.
		log.Error().Err(err).Msg("user not found")
		return nil, e.NewNotFoundError("user not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get user bio from the database")
		return nil, e.NewInternalServerError()
	}

	return ub, nil
}

func (s *BaseUserService) GetCurrentEmail(
	ctx context.Context,
	userID string,
) (string, e.Error) {
	log.Info().Msg("getting user from the database")
	u, err := s.ur.GetUserByID(ctx, userID)
	if _, ok := err.(repository.NotFoundError); ok {
		// this shouldn't happen in an ideal scenario
		log.Error().Err(err).Msg("user not found")
		return "", e.NewNotFoundError("user not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get user from the database")
		return "", e.NewInternalServerError()
	}

	return u.Email, nil
}

func (s *BaseUserService) UpdateBasicInfo(
	ctx context.Context,
	userID string,
	ub *repository.UserBio,
) e.Error {
	log.Info().Msg("updating user bio in database")
	_, err := s.ur.UpdateUserBioByID(ctx, userID, ub)
	if err != nil {
		log.Error().Err(err).Msg("failed to update user bio in the database")
		return e.NewInternalServerError()
	}

	return nil
}

func (s *BaseUserService) RequestEmailChange(
	ctx context.Context,
	userID string,
	newEmail string,
) e.Error {
	log.Info().Msg("getting user from database")
	u, err := s.ur.GetUserByID(ctx, userID)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("user not found")
		return e.NewNotFoundError("user not found")
	}

	log.Info().Msg("checking new email with the old one")
	if u.Email == newEmail {
		log.Error().Msg("new email is the same as the old one")
		return e.NewBadRequestError("new email can't be the same as the old one")
	}

	log.Info().Msg("checking if new email is available/not taken")
	_, err = s.ur.GetUserByEmail(ctx, newEmail)
	if _, ok := err.(repository.NotFoundError); !ok {
		log.Error().Err(err).Msg("new email is already taken")
		return e.NewBadRequestError("email is already registered")
	}

	log.Info().Msg("generating email change token")
	token := security.GenerateRandomID()

	log.Info().Msg("storing token in redis")
	t := &repository.EmailChangeToken{
		NewEmail: newEmail,
		Token:    string(token),
	}
	err = s.tr.CreateEmailChangeToken(ctx, userID, t)
	if err != nil {
		log.Error().Err(err).Msg("failed to store token in redis")
		return e.NewInternalServerError()
	}

	// TODO: send the verification link to the new email
	log.Debug().Msgf("email change token: %s", token)

	return nil
}

func (s *BaseUserService) VerifyEmailChange(
	ctx context.Context,
	userID string,
	token string,
) e.Error {
	log.Info().Msg("retrieving token from redis")
	t, err := s.tr.GetEmailChangeToken(ctx, userID)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("token not found")
		return e.NewNotFoundError("token not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve token from redis")
		return e.NewInternalServerError()
	}

	log.Info().Msg("checking token")
	if t.Token != token {
		log.Error().Msg("token is invalid")
		return e.NewBadRequestError("token is invalid")
	}

	log.Info().Msg("updating user email")
	_, err = s.ur.UpdateEmailByID(ctx, userID, t.NewEmail)
	if err != nil {
		log.Error().Err(err).Msg("failed to update user email")
		return e.NewInternalServerError()
	}

	log.Info().Msg("deleting token from redis")
	err = s.tr.DeleteEmailChangeToken(ctx, userID)
	if _, ok := err.(repository.NotFoundError); !ok && err != nil {
		log.Error().Err(err).Msg("failed to delete token from redis")
		return e.NewInternalServerError()
	}

	return nil
}

func (s *BaseUserService) ChangePassword(
	ctx context.Context,
	userID string,
	currentPassword string,
	newPassword string,
) e.Error {
	log.Info().Msg("getting user from database")
	u, err := s.ur.GetUserByID(ctx, userID)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("user not found")
		return e.NewNotFoundError("user not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get user from the database")
		return e.NewInternalServerError()
	}

	log.Info().Msg("checking current password")
	err = security.CheckPassword(currentPassword, u.Password)
	if err != nil {
		log.Error().Err(err).Msg("wrong password")
		return e.NewUnauthorizedError("wrong password")
	}

	log.Info().Msg("retrieving last 3 passwords")
	hashes, err := s.phr.GetLastNPasswordHashes(ctx, userID, 3)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve last 3 passwords")
		return e.NewInternalServerError()
	}

	log.Info().Msg("checking last 3 passwords")
	if slice.AnyStr(hashes, func(h string) bool {
		err := security.CheckPassword(newPassword, h)
		return err == nil
	}) {
		log.Error().Msg("new password is the same as one of the last 3 passwords")
		return e.NewBadRequestError("new password can't be the same as one of the last 3 passwords")
	}

	log.Info().Msg("hashing password")
	hash, err := security.HashPassword(newPassword)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return e.NewInternalServerError()
	}

	log.Info().Msg("updating user password")
	_, err = s.ur.UpdatePasswordByID(ctx, userID, hash)
	if err != nil {
		log.Error().Err(err).Msg("failed to update user password")
		return e.NewInternalServerError()
	}

	ph := &repository.PasswordHistory{
		UserID:   u.ID,
		Password: hash,
	}

	log.Info().Msg("creating password history record")
	_, err = s.phr.CreatePasswordHistoryRecord(ctx, ph)
	if err != nil {
		log.Error().Err(err).Msg("failed to create password history record")
		return e.NewInternalServerError()
	}

	// TODO: not in the requirement, but it'll be nice to invalidate all other
	// sessions after changing the password

	return nil
}
