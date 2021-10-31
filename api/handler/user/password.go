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

type changePasswordRequest struct {
	PasswordCurrent string `json:"password_current"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type changePasswordResponse struct {
	Success bool `json:"success"`
}

func validateChangePasswordRequest(req *changePasswordRequest) (map[string]string, bool) {
	fields := map[string]string{}

	errMsg, ok := validator.ValidatePasswordSimple(req.PasswordCurrent, "password_current")
	if !ok {
		fields["password_current"] = errMsg
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

func ChangePassword(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &changePasswordRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.Error(w, e.NewBadRequestError("cannot decode request body")).JSON()
			return
		}

		fields, ok := validateChangePasswordRequest(req)
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

		err = us.ChangePassword(ctx, at.UserID, req.PasswordCurrent, req.Password)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		response.OK(w, &changePasswordResponse{
			Success: true,
		}).JSON()
	}
}
