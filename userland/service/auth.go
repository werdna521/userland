package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/mailer"
	"github.com/werdna521/userland/producer"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/postgres"
	"github.com/werdna521/userland/repository/redis"
	"github.com/werdna521/userland/security"
	"github.com/werdna521/userland/security/jwt"
	"github.com/werdna521/userland/utils/slice"
)

type AuthService interface {
	Register(ctx context.Context, user *repository.User) e.Error
	SendEmailVerification(ctx context.Context, email string) e.Error
	VerifyEmail(ctx context.Context, email string, token string) e.Error
	Login(
		ctx context.Context,
		user *repository.User,
		clientID string,
		ip string,
	) (*jwt.AccessToken, e.Error)
	ForgotPassword(ctx context.Context, email string) e.Error
	ResetPassword(ctx context.Context, token string, newPassword string) e.Error
}

type BaseAuthService struct {
	ur  postgres.UserRepository
	phr postgres.PasswordHistoryRepository
	tr  redis.TokenRepository
	sr  redis.SessionRepository
	m   mailer.Mailer
	lp  producer.LogProducer
}

func NewBaseAuthService(
	ur postgres.UserRepository,
	phr postgres.PasswordHistoryRepository,
	tr redis.TokenRepository,
	sr redis.SessionRepository,
	m mailer.Mailer,
	lp producer.LogProducer,
) *BaseAuthService {
	return &BaseAuthService{
		ur:  ur,
		phr: phr,
		tr:  tr,
		sr:  sr,
		m:   m,
		lp:  lp,
	}
}

func (s *BaseAuthService) Register(ctx context.Context, u *repository.User) e.Error {
	log.Info().Msg("retrieving user from database")
	_, err := s.ur.GetUserByEmail(ctx, u.Email)
	if _, ok := err.(repository.NotFoundError); !ok && err != nil {
		log.Error().Err(err).Msg("failed to retrieve user from database")
		return e.NewInternalServerError()
	}
	if err == nil {
		log.Info().Msg("user already exists")
		return e.NewConflictError("user already exists")
	}

	log.Info().Msg("hashing password")
	hash, err := security.HashPassword(u.Password)
	if err != nil {
		log.Error().Err(err).Msg("fail to hash password")
		return e.NewInternalServerError()
	}
	u.Password = hash

	u.IsActive = false

	log.Info().Msg("creating and registering user")
	u, err = s.ur.CreateUser(ctx, u)
	if _, ok := err.(repository.UniqueViolationError); ok {
		log.Error().Err(err).Msg("user already exists")
		return e.NewConflictError("user already exists")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return e.NewInternalServerError()
	}

	ph := &repository.PasswordHistory{
		UserID:   u.ID,
		Password: hash,
	}

	log.Info().Msg("adding password to history")
	s.phr.CreatePasswordHistoryRecord(ctx, ph)

	log.Info().Msg("generating verification token")
	verificationToken := security.GenerateRandomID()

	log.Info().Msg("storing email verification token")
	err = s.tr.CreateEmailVerificationToken(ctx, u.ID, string(verificationToken))
	if err != nil {
		log.Error().Err(err).Msg("failed to store email verification token")
		return e.NewInternalServerError()
	}

	verificationLink := fmt.Sprintf(
		"http://localhost:3000/api/v1/auth/verification?id=%s&token=%s",
		u.ID,
		verificationToken,
	)

	log.Debug().Msgf("verification link: %s", verificationLink)
	log.Info().Msg("sending verification link")
	email := mailer.Email{
		Name:  u.UserBio.Fullname,
		Email: u.Email,
	}
	err = mailer.SendEmailVerificationMail(ctx, s.m, email, verificationLink)
	if err != nil {
		log.Error().Err(err).Msg("failed to send verification link")
		return e.NewInternalServerError()
	}

	return nil
}

func (s *BaseAuthService) SendEmailVerification(
	ctx context.Context,
	email string,
) e.Error {
	log.Info().Msg("retrieving user from database")
	u, err := s.ur.GetUserByEmail(ctx, email)
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
	err = s.tr.CreateEmailVerificationToken(ctx, u.ID, string(verificationToken))
	if err != nil {
		log.Error().Err(err).Msg("failed to store email verification token")
		return e.NewInternalServerError()
	}

	verificationLink := fmt.Sprintf(
		"http://localhost:3000/api/v1/auth/verification?id=%s&token=%s",
		u.ID,
		verificationToken,
	)

	log.Debug().Msgf("verification link: %s", verificationLink)
	log.Info().Msg("sending verification link")
	em := mailer.Email{
		Name:  u.Email,
		Email: u.Email,
	}
	err = mailer.SendEmailVerificationMail(ctx, s.m, em, verificationLink)
	if err != nil {
		log.Error().Err(err).Msg("failed to send verification link")
		return e.NewInternalServerError()
	}

	return nil
}

func (s *BaseAuthService) VerifyEmail(
	ctx context.Context,
	userID string,
	verificationToken string,
) e.Error {
	log.Info().Msg("retrieving verification token from repository")
	storedToken, err := s.tr.GetEmailVerificationToken(ctx, userID)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("verification token not found")
		return e.NewNotFoundError("invalid token")
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
	_, err = s.ur.UpdateUserActivationStatusByID(ctx, userID, true)
	if err != nil {
		log.Error().Err(err).Msg("failed to activate user account")
		return e.NewInternalServerError()
	}

	// still, remove the data from redis, even though it doesn't really matter tbh
	log.Info().Msg("removing verification details from redis")
	err = s.tr.DeleteEmailVerificationToken(ctx, userID)
	if _, ok := err.(repository.NotFoundError); !ok && err != nil {
		log.Error().Err(err).Msg("failed to remove verification details from redis")
		return e.NewInternalServerError()
	}

	return nil
}

func (s *BaseAuthService) Login(
	ctx context.Context,
	u *repository.User,
	clientID string,
	ip string,
) (*jwt.AccessToken, e.Error) {
	log.Info().Msg("retrieving user from database")
	userFromDB, err := s.ur.GetUserByEmail(ctx, u.Email)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("user not found")
		return nil, e.NewNotFoundError("user not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("checking if user is active")
	if !userFromDB.IsActive {
		log.Error().Msg("user is not active")
		return nil, e.NewForbiddenError("user is not active")
	}

	log.Info().Msg("checking if password is correct")
	err = security.CheckPassword(u.Password, userFromDB.Password)
	if err != nil {
		log.Error().Err(err).Msg("password is incorrect")
		return nil, e.NewUnauthorizedError("password is incorrect")
	}

	log.Info().Msg("generating session ID")
	sessionID := security.GenerateRandomID()

	log.Info().Msg("generating access token")
	at, err := jwt.CreateAccessToken(userFromDB.ID, string(sessionID))
	if err != nil {
		log.Error().Err(err).Msg("failed to generate access token")
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("storing access token in redis")
	token := &repository.AccessToken{
		ID:        at.JTI,
		UserID:    at.UserID,
		SessionID: at.SessionID,
	}
	err = s.sr.CreateAccessToken(ctx, token, jwt.AccessTokenLife)
	if err != nil {
		log.Error().Err(err).Msg("failed to store access token")
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("storing session in redis")
	session := &repository.Session{
		ID:        at.SessionID,
		IPAddress: ip,
		Client:    clientID,
		UserID:    at.UserID,
	}
	err = s.sr.CreateSession(ctx, session, jwt.AccessTokenLife)
	if err != nil {
		log.Error().Err(err).Msg("failed to store session in redis")
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("adding the session id to a user session index set")
	err = s.sr.AddUserSessionToIndex(ctx, session)
	if err != nil {
		log.Error().Err(err).Msg("failed to add the session id to the index set")
		return nil, e.NewInternalServerError()
	}

	log.Info().Msg("producing message to log micro")
	ll := &producer.LoginLog{
		UserID:   userFromDB.ID,
		RemoteIP: ip,
	}
	err = s.lp.ProduceLoginTopic(ll)
	if err != nil {
		log.Error().Err(err).Msg("failed to produce message to log micro")
		return nil, e.NewInternalServerError()
	}

	return at, nil
}

func (s *BaseAuthService) ForgotPassword(ctx context.Context, email string) e.Error {
	log.Info().Msg("retrieving user from the db")
	u, err := s.ur.GetUserByEmail(ctx, email)
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
	err = s.tr.CreateForgotPasswordToken(ctx, u.ID, string(token))
	if err != nil {
		log.Error().Err(err).Msg("failed to store forgot password token")
		return e.NewInternalServerError()
	}

	log.Debug().Msgf("token: %s", token)
	log.Info().Msg("sending password reset mail")
	em := mailer.Email{
		Name:  email,
		Email: email,
	}
	err = mailer.SendPasswordResetMail(ctx, s.m, em, string(token))
	if err != nil {
		log.Error().Err(err).Msg("failed to send password reset mail")
		return e.NewInternalServerError()
	}

	return nil
}

func (s *BaseAuthService) ResetPassword(
	ctx context.Context,
	token string,
	newPassword string,
) e.Error {
	userID, err := s.tr.GetForgotPasswordToken(ctx, token)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("token not found")
		return e.NewUnauthorizedError("invalid token")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve forgot password token details")
		return e.NewInternalServerError()
	}

	log.Info().Msg("retrieving last 3 password hash from db")
	hashes, err := s.phr.GetLastNPasswordHashes(ctx, userID, 3)
	if err != nil {
		log.Error().Err(err).Msg("failed to get the password hashes")
		return e.NewInternalServerError()
	}

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

	log.Info().Msg("retrieving old password from database")
	u, err := s.ur.GetUserByID(ctx, userID)
	if _, ok := err.(repository.NotFoundError); ok {
		log.Error().Err(err).Msg("user is no longer in the database")
		return e.NewUnauthorizedError("invalid token")
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve user from the db")
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

	log.Info().Msg("creating forgot password record")
	_, err = s.phr.CreatePasswordHistoryRecord(ctx, ph)
	if err != nil {
		log.Error().Err(err).Msg("failed to create forgot password record")
		return e.NewInternalServerError()
	}

	log.Info().Msg("removing forgot password token")
	err = s.tr.DeleteForgotPasswordToken(ctx, token)
	if _, ok := err.(repository.NotFoundError); !ok && err != nil {
		log.Error().Err(err).Msg("failed to remove forgot password token")
		return e.NewInternalServerError()
	}

	return nil
}
