package response

import (
	"encoding/json"
	"net/http"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/error/client"
)

type baseErrorResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message"`
}

type unprocessableEntityResponse struct {
	Success bool              `json:"success"`
	Fields  map[string]string `json:"fields"`
}

func respondWithError(w http.ResponseWriter, err e.Error, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode())
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, err e.Error) {
	switch err := err.(type) {
	case client.UnprocessableEntityError:
		respondWithError(
			w,
			err,
			unprocessableEntityResponse{
				Success: false,
				Fields:  err.Fields,
			},
		)
	default:
		respondWithError(
			w,
			err,
			baseErrorResponse{
				Success: false,
				Msg:     err.Error(),
			},
		)
	}
}
