package user

import (
	"encoding/json"
	"net/http"
	"time"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/request"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/api/validator"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/service"
)

type getInfoDetailResponse struct {
	Success   bool      `json:"success"`
	ID        string    `json:"id"`
	Fullname  string    `json:"fullname"`
	Location  string    `json:"location"`
	Bio       string    `json:"bio"`
	Web       string    `json:"web"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
}

func GetInfoDetail(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		at, err := request.GetAccessTokenFromCtx(ctx)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		ub, err := us.GetInfoDetail(ctx, at.UserID)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		response.OK(w, &getInfoDetailResponse{
			Success:   true,
			ID:        at.UserID,
			Fullname:  ub.Fullname,
			Location:  ub.Location,
			Bio:       ub.Bio,
			Web:       ub.Web,
			Picture:   ub.Picture,
			CreatedAt: ub.CreatedAt,
		}).JSON()
	}
}

type updateBasicInfoRequest struct {
	Fullname string `json:"fullname"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
	Web      string `json:"web"`
}

type updateBasicInfoResponse struct {
	Success bool `json:"success"`
}

func validateUpdateBasicInfoRequest(req *updateBasicInfoRequest) (map[string]string, bool) {
	fields := map[string]string{}

	errMsg, ok := validator.ValidateFullname(req.Fullname)
	if !ok {
		fields["fullname"] = errMsg
	}

	errMsg, ok = validator.ValidateLocation(req.Location)
	if !ok {
		fields["location"] = errMsg
	}

	errMsg, ok = validator.ValidateBio(req.Bio)
	if !ok {
		fields["bio"] = errMsg
	}

	errMsg, ok = validator.ValidateWeb(req.Web)
	if !ok {
		fields["web"] = errMsg
	}

	return fields, len(fields) == 0
}

func UpdateBasicInfo(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &updateBasicInfoRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.Error(w, e.NewBadRequestError("cannot decode request body")).JSON()
			return
		}

		fields, ok := validateUpdateBasicInfoRequest(req)
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

		ub := &repository.UserBio{
			Fullname: req.Fullname,
			Location: req.Location,
			Bio:      req.Bio,
			Web:      req.Web,
		}
		err = us.UpdateBasicInfo(ctx, at.UserID, ub)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		response.OK(w, &updateBasicInfoResponse{
			Success: true,
		}).JSON()
	}
}
