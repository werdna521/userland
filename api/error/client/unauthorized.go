package client

import "net/http"

type UnauthorizedError struct {
	Msg string
}

func (e UnauthorizedError) Error() string {
	return e.Msg
}

func (e UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}
