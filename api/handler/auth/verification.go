package auth

import (
	"net/http"
	"net/url"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/service"
)

type verifyEmailRequest struct {
	Email string
	Code  string
}

type verifyEmailResponse struct {
	Success bool `json:"success"`
}

func toVerifyEmailRequest(params url.Values) *verifyEmailRequest {
	return &verifyEmailRequest{
		Email: params.Get("email"),
		Code:  params.Get("code"),
	}
}

func validateVerifyEmailRequest(req *verifyEmailRequest) bool {
	return req.Email != "" && req.Code != ""
}

func VerifyEmail(as service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := toVerifyEmailRequest(r.URL.Query())

		ok := validateVerifyEmailRequest(req)
		if !ok {
			response.Error(w, e.NewBadRequestError("bad request")).JSON()
			return
		}

		ctx := r.Context()
		err := as.VerifyEmail(ctx, req.Email, req.Code)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		response.OK(w, &verifyEmailResponse{
			Success: true,
		}).JSON()
	}
}
