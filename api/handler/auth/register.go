package auth

import (
	"encoding/json"
	"net/http"

	"github.com/werdna521/userland/api/error/client"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/api/validator"
	"github.com/werdna521/userland/repository"
)

type registerRequest struct {
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func validateRegisterRequest(req *registerRequest) (map[string]string, bool) {
	fields := map[string]string{}

	errMsg, ok := validator.ValidateFullname(req.Fullname)
	if !ok {
		fields["fullname"] = errMsg
	}

	return fields, true
}

func Register(ur repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &registerRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.Error(w, client.NewBadRequestError("cannot decode request body"))
			return
		}

		fields, ok := validateRegisterRequest(req)
		if !ok {
			response.Error(w, client.NewUnprocessableEntityError(fields))
			return
		}

		u := &repository.User{
			Fullname: req.Fullname,
			Email:    req.Email,
			Password: req.Password,
		}

		ctx := r.Context()
		err = ur.CreateUser(ctx, u)
		if err != nil {
			response.Error(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}
