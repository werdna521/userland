package client

import "net/http"

type ConflictError struct {
	Msg string
}

func (e ConflictError) Error() string {
	return e.Msg
}

func (e ConflictError) StatusCode() int {
	return http.StatusConflict
}
