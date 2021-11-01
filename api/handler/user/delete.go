package user

import (
	"encoding/json"
	"net/http"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/request"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/api/validator"
	"github.com/werdna521/userland/service"
)

type deleteAccountRequest struct {
	Password string `json:"password"`
}

type deleteAccountResponse struct {
	Success bool `json:"success"`
}

func validateDeleteAccountRequest(req *deleteAccountRequest) (map[string]string, bool) {
	fields := map[string]string{}

	errMsg, ok := validator.ValidatePasswordSimple(req.Password, "password")
	if !ok {
		fields["password"] = errMsg
	}

	return fields, len(fields) == 0
}

func DeleteAccount(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &deleteAccountRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.Error(w, e.NewBadRequestError("cannot decode request body")).JSON()
			return
		}

		fields, ok := validateDeleteAccountRequest(req)
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

		err = us.DeleteAccount(ctx, at.UserID, req.Password)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		response.OK(w, &deleteAccountResponse{
			Success: true,
		}).JSON()
	}
}
