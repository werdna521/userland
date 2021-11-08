package response

import (
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

func respondWithError(w http.ResponseWriter, err e.Error, v interface{}) httpResponse {
	return httpResponse{
		statusCode: err.StatusCode(),
		w:          w,
		v:          v,
	}
}

func Error(w http.ResponseWriter, err e.Error) httpResponse {
	switch err := err.(type) {
	case client.UnprocessableEntityError:
		return respondWithError(
			w,
			err,
			unprocessableEntityResponse{
				Success: false,
				Fields:  err.Fields,
			},
		)
	default:
		return respondWithError(
			w,
			err,
			baseErrorResponse{
				Success: false,
				Msg:     err.Error(),
			},
		)
	}
}
