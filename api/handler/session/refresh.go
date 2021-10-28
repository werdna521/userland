package session

import (
	"net/http"

	"github.com/werdna521/userland/api/request"
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
		ctx := r.Context()
		at, err := request.GetAccessTokenFromCtx(ctx)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

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
