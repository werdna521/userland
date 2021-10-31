package user

import (
	"encoding/json"
	"net/http"
	"net/url"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/request"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/api/validator"
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

type requestEmailAddressChangeRequest struct {
	Email string `json:"email"`
}

type requestEmailAddressChangeResponse struct {
	Success bool `json:"success"`
}

func validateRequestEmailAddressChangeRequest(
	req *requestEmailAddressChangeRequest,
) (map[string]string, bool) {
	fields := map[string]string{}

	errMsg, ok := validator.ValidateEmail(req.Email)
	if !ok {
		fields["email"] = errMsg
	}

	return fields, len(fields) == 0
}

func RequestEmailAddressChange(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &requestEmailAddressChangeRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.Error(w, e.NewBadRequestError("cannot decode request body")).JSON()
			return
		}

		fields, ok := validateRequestEmailAddressChangeRequest(req)
		if !ok {
			response.Error(w, e.NewUnprocessableEntityError(fields)).JSON()
			return
		}

		ctx := r.Context()
		at, err := request.GetAccessTokenFromCtx(ctx)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		err = us.RequestEmailChange(ctx, at.UserID, req.Email)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		response.OK(w, &requestEmailAddressChangeResponse{
			Success: true,
		}).JSON()
	}
}

type verifyEmailChangeRequest struct {
	UserID string
	Token  string
}

type verifyEmailChangeResponse struct {
	Success bool `json:"success"`
}

func toVerifyEmailChangeRequest(params url.Values) *verifyEmailChangeRequest {
	return &verifyEmailChangeRequest{
		UserID: params.Get("id"),
		Token:  params.Get("token"),
	}
}

func validateVerifyEmailChangeRequest(req *verifyEmailChangeRequest) bool {
	return req.UserID != "" && req.Token != ""
}

func VerifyEmailChange(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := toVerifyEmailChangeRequest(r.URL.Query())

		ok := validateVerifyEmailChangeRequest(req)
		if !ok {
			response.Error(w, e.NewBadRequestError("bad request")).JSON()
			return
		}

		ctx := r.Context()
		err := us.VerifyEmailChange(ctx, req.UserID, req.Token)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		response.OK(w, &verifyEmailChangeResponse{
			Success: true,
		}).JSON()
	}
}
