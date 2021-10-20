package auth

import (
	"encoding/json"
	"net/http"

	"github.com/werdna521/userland/repository"
)

func HandleRegister(ur repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u := &repository.User{
			FullName: "Andrew Cen",
			Email:    "andrew@me.com",
			Password: "password",
		}

		err := ur.CreateUser(ctx, u)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"success": "false",
				"error":   err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}
