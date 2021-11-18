package auth

import (
	"encoding/json"
	"net/http"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/request"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/api/validator"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/security/jwt"
	"github.com/werdna521/userland/service"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Success     bool             `json:"success"`
	RequireTFA  bool             `json:"require_tfa"`
	AccessToken *jwt.AccessToken `json:"access_token"`
}

func validateLoginRequest(req *loginRequest) (map[string]string, bool) {
	fields := map[string]string{}

	errMsg, ok := validator.ValidateEmail(req.Email)
	if !ok {
		fields["email"] = errMsg
	}

	errMsg, ok = validator.ValidatePasswordSimple(req.Password, "password")
	if !ok {
		fields["password"] = errMsg
	}

	return fields, len(fields) == 0
}

func Login(as service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Header.Get("X-API-ClientID")
		ip := request.ParseIPAddress(r)

		req := &loginRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.Error(w, e.NewBadRequestError("cannot decode request body")).JSON()
			return
		}

		fields, ok := validateLoginRequest(req)
		if !ok {
			response.Error(w, e.NewUnprocessableEntityError(fields)).JSON()
			return
		}

		ctx := r.Context()
		u := &repository.User{
			Email:    req.Email,
			Password: req.Password,
		}
		at, err := as.Login(ctx, u, clientID, ip)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		response.OK(w, &loginResponse{
			Success: true,
			// TODO: implement tfa properly after everything else is done :)
			RequireTFA:  false,
			AccessToken: at,
		}).JSON()
	}
}
