package session

import (
	"net/http"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/middleware"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/security/jwt"
	"github.com/werdna521/userland/service"
)

type generateAccessTokenResponse struct {
	Success     bool             `json:"success"`
	AccessToken *jwt.AccessToken `json:"accessToken"`
}

func GenerateAccessToken(ss service.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rt, ok := ctx.Value(middleware.RefreshTokenCtxKey).(*jwt.RefreshToken)
		if !ok {
			response.Error(w, e.NewBadRequestError("cannot parse refresh token"))
			return
		}

		at, err := ss.GenerateAccessToken(ctx, rt)
		if err != nil {
			response.Error(w, err)
			return
		}

		response.OK(w, &generateAccessTokenResponse{
			Success:     true,
			AccessToken: at,
		}).JSON()
	}
}
