package auth

import (
	"encoding/json"
	"net/http"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/api/validator"
	"github.com/werdna521/userland/service"
)

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type forgotPasswordResponse struct {
	Success bool `json:"success"`
}

func validateForgotPasswordRequest(req *forgotPasswordRequest) (map[string]string, bool) {
	fields := map[string]string{}

	errMsg, ok := validator.ValidateEmail(req.Email)
	if !ok {
		fields["email"] = errMsg
	}

	return fields, len(fields) == 0
}

func ForgotPassword(au service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &forgotPasswordRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.Error(w, e.NewBadRequestError("cannot decode request body")).JSON()
			return
		}

		fields, ok := validateForgotPasswordRequest(req)
		if !ok {
			response.Error(w, e.NewUnprocessableEntityError(fields)).JSON()
			return
		}

		ctx := r.Context()
		err = au.ForgotPassword(ctx, req.Email)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		response.OK(w, &forgotPasswordResponse{
			Success: true,
		}).JSON()
	}
}
