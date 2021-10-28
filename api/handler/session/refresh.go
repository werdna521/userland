package session

import (
	"net/http"

	"github.com/werdna521/userland/api/middleware"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/security/jwt"
	"github.com/werdna521/userland/service"
)

type generateRefreshTokenResponse struct {
	Success      bool              `json:"success"`
	RefreshToken *jwt.RefreshToken `json:"refresh_token"`
}

func GenerateRefreshToken(ss service.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		at := r.Context().Value(middleware.AccessTokenCtxKey).(*jwt.AccessToken)

		ctx := r.Context()
		rt, err := ss.GenerateRefreshToken(ctx, at)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		response.OK(w, &generateRefreshTokenResponse{
			Success:      true,
			RefreshToken: rt,
		}).JSON()
	}
}
