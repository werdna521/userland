package session

import (
	"net/http"

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
		response.OK(w, &generateAccessTokenResponse{
			Success: true,
		}).JSON()
	}
}
