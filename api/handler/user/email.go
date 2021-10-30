package user

import (
	"net/http"

	"github.com/werdna521/userland/api/request"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/service"
)

type getCurrentEmailAddressResponse struct {
	Success bool   `json:"success"`
	Email   string `json:"email"`
}

func GetCurrentEmailAddress(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		at, err := request.GetAccessTokenFromCtx(ctx)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		email, err := us.GetCurrentEmail(ctx, at.UserID)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		response.OK(w, &getCurrentEmailAddressResponse{
			Success: true,
			Email:   email,
		}).JSON()
	}
}
