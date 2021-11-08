package client

import "net/http"

type ForbiddenError struct {
	Msg string
}

func (e ForbiddenError) Error() string {
	return e.Msg
}

func (e ForbiddenError) StatusCode() int {
	return http.StatusForbidden
}
