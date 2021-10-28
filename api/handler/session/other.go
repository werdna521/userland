package session

import (
	"net/http"

	"github.com/werdna521/userland/api/request"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/service"
)

type deleteAllOtherSessionsResponse struct {
	Success bool `json:"success"`
}

func DeleteAllOtherSessions(ss service.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		at, err := request.GetAccessTokenFromCtx(ctx)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		session := &repository.Session{
			ID:     at.SessionID,
			UserID: at.UserID,
		}
		err = ss.RemoveAllOtherSessions(ctx, session)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		response.OK(w, &deleteAllOtherSessionsResponse{
			Success: true,
		}).JSON()
	}
}
