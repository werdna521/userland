package service

import (
	"context"

	"github.com/rs/zerolog/log"
	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/postgres"
)

type UserService interface {
	GetInfoDetail(ctx context.Context, userID string) (*repository.UserBio, e.Error)
}

type BaseUserService struct {
	ur postgres.UserRepository
}

func NewBaseUserService(ur postgres.UserRepository) *BaseUserService {
	return &BaseUserService{
		ur: ur,
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
