package auth

import (
	"encoding/json"
	"net/http"
	"net/url"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/api/validator"
	"github.com/werdna521/userland/service"
)

type sendVerificationRequest struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
}

type sendVerificationResponse struct {
	Success bool `json:"success"`
}

func validateSendVerificationRequest(req *sendVerificationRequest) (map[string]string, bool) {
	fields := map[string]string{}

	errMsg, ok := validator.ValidateVerificationType(req.Type)
	if !ok {
		fields["type"] = errMsg
	}

	errMsg, ok = validator.ValidateRecipient(req.Recipient)
	if !ok {
		fields["recipient"] = errMsg
	}

	return fields, len(fields) == 0
}

func SendVerification(as service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &sendVerificationRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.Error(w, e.NewBadRequestError("cannot decode request body")).JSON()
			return
		}

		fields, ok := validateSendVerificationRequest(req)
		if !ok {
			response.Error(w, e.NewUnprocessableEntityError(fields)).JSON()
			return
		}

		switch req.Type {
		case "email.verify":
			ctx := r.Context()
			err = as.SendEmailVerification(ctx, req.Recipient)
			if err != nil {
				response.Error(w, err.(e.Error)).JSON()
				return
			}

			response.OK(w, &sendVerificationResponse{
				Success: true,
			}).JSON()
			return
		default:
			response.Error(w, e.NewBadRequestError("invalid type")).JSON()
		}
	}
}

type verifyEmailRequest struct {
	UserID string
	Token  string
}

type verifyEmailResponse struct {
	Success bool `json:"success"`
}

func toVerifyEmailRequest(params url.Values) *verifyEmailRequest {
	return &verifyEmailRequest{
		UserID: params.Get("id"),
		Token:  params.Get("token"),
	}
}

func validateVerifyEmailRequest(req *verifyEmailRequest) bool {
	return req.UserID != "" && req.Token != ""
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
		err := as.VerifyEmail(ctx, req.UserID, req.Token)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		response.OK(w, &verifyEmailResponse{
			Success: true,
		}).JSON()
	}
}
