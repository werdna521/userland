package session

import (
	"net/http"
	"time"

	"github.com/werdna521/userland/api/request"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/service"
)

type listSessionsResponse struct {
	Success  bool           `json:"success"`
	Sessions []*userSession `json:"sessions"`
}

type userSession struct {
	IsCurrent bool      `json:"isCurrent"`
	IPAddress string    `json:"ip"`
	Client    *client   `json:"client"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ListSessions(ss service.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		at, err := request.GetAccessTokenFromCtx(ctx)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		sessions, err := ss.ListSessions(ctx, at)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		userSessions := []*userSession{}
		for _, s := range sessions {
			us := &userSession{
				IsCurrent: s.ID == at.SessionID,
				IPAddress: s.IPAddress,
				Client: &client{
					ID:   s.ID,
					Name: s.Client,
				},
				CreatedAt: s.CreatedAt,
				UpdatedAt: s.UpdatedAt,
			}
			userSessions = append(userSessions, us)
		}

		response.OK(w, &listSessionsResponse{
			Success:  true,
			Sessions: userSessions,
		}).JSON()
	}
}

type deleteSessionResponse struct {
	Success bool `json:"success"`
}

func EndCurrentSession(ss service.SessionService) http.HandlerFunc {
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
		err = ss.RemoveSession(ctx, session)
		if err != nil {
			response.Error(w, err).JSON()
			return
		}

		response.OK(w, &deleteSessionResponse{
			Success: true,
		}).JSON()
	}
}
