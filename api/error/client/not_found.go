package client

import "net/http"

type NotFoundError struct {
	Msg string
}

func (e NotFoundError) Error() string {
	return e.Msg
}

func (e NotFoundError) StatusCode() int {
	return http.StatusNotFound
}
