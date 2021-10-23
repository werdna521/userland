package auth

import (
	"encoding/json"
	"net/http"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/api/validator"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/service"
)

type registerRequest struct {
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type registerResponse struct {
	Success bool `json:"success"`
}

func validateRegisterRequest(req *registerRequest) (map[string]string, bool) {
	fields := map[string]string{}

	errMsg, ok := validator.ValidateFullname(req.Fullname)
	if !ok {
		fields["fullname"] = errMsg
	}

	errMsg, ok = validator.ValidateEmail(req.Email)
	if !ok {
		fields["email"] = errMsg
	}

	errMsg, ok = validator.ValidatePassword(req.Password)
	if !ok {
		fields["password"] = errMsg
	}

	errMsg, ok = validator.ValidatePasswordConfirm(req.Password, req.PasswordConfirm)
	if !ok {
		fields["password_confirm"] = errMsg
	}

	return fields, len(fields) == 0
}

func Register(as service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &registerRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.Error(w, e.NewBadRequestError("cannot decode request body"))
			return
		}

		fields, ok := validateRegisterRequest(req)
		if !ok {
			response.Error(w, e.NewUnprocessableEntityError(fields)).JSON()
			return
		}

		u := &repository.User{
			Fullname: req.Fullname,
			Email:    req.Email,
			Password: req.Password,
		}

		ctx := r.Context()
		err = as.Register(ctx, u)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		response.OK(w, &registerResponse{
			Success: true,
		}).JSON()
	}
}
