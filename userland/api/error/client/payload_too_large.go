package client

import "net/http"

type RequestEntityTooLargeError struct {
	Msg string
}

func (e RequestEntityTooLargeError) Error() string {
	return e.Msg
}

func (e RequestEntityTooLargeError) StatusCode() int {
	return http.StatusRequestEntityTooLarge
}
