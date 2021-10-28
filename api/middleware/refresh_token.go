package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/redis"
	"github.com/werdna521/userland/security/jwt"
)

type RefreshTokenKey string

const RefreshTokenCtxKey RefreshTokenKey = "refreshtoken"

func ValidateRefreshToken(sr redis.SessionRepository) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				log.Error().Msg("No authorization header")
				response.Error(w, e.NewUnauthorizedError("no token provided")).JSON()
				return
			}

			bearer := strings.Split(authHeader, " ")
			if len(bearer) != 2 {
				log.Error().Msg("Invalid authorization header")
				response.Error(w, e.NewBadRequestError("bad authorization header format")).JSON()
				return
			}

			log.Info().Msg("parsing access token")
			jwtString := bearer[1]
			rt, isValid, err := jwt.ParseRefreshToken(jwtString)
			if !isValid {
				log.Error().Msg("Invalid access token")
				response.Error(w, e.NewUnauthorizedError("invalid token")).JSON()
				return
			}
			if err != nil {
				log.Error().Err(err).Msg("failed to parse token")
				response.Error(w, e.NewInternalServerError()).JSON()
				return
			}

			log.Info().Msg("checking if token is valid")
			ctx := r.Context()
			tokenExists, err := sr.CheckRefreshToken(ctx, &repository.RefreshToken{
				ID:        rt.JTI,
				SessionID: rt.SessionID,
				UserID:    rt.UserID,
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to retrieve token from redis")
				response.Error(w, e.NewInternalServerError()).JSON()
				return
			}

			if !tokenExists {
				log.Error().Msg("token does not exist")
				response.Error(w, e.NewUnauthorizedError("invalid token")).JSON()
				return
			}

			ctx = context.WithValue(r.Context(), RefreshTokenCtxKey, rt)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
