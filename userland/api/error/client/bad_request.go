package client

import "net/http"

type BadRequestError struct {
	Msg string
}

func (e BadRequestError) Error() string {
	return e.Msg
}

func (e BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}
