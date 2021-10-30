package user

import (
	"net/http"
	"time"

	"github.com/werdna521/userland/api/request"
	"github.com/werdna521/userland/api/response"
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
