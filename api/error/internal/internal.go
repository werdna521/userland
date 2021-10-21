package internal

import (
	"net/http"
)

type InternalServerError struct{}

func (e InternalServerError) Error() string {
	return "Internal Server Error"
}

func (e InternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}
