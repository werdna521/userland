package client

import "net/http"

type BadRequestError struct {
	Msg string
}

func NewBadRequestError(msg string) BadRequestError {
	return BadRequestError{
		Msg: msg,
	}
}

func (e BadRequestError) Error() string {
	return e.Msg
}

func (e BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}
