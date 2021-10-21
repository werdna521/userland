package error

import (
	"github.com/werdna521/userland/api/error/client"
	"github.com/werdna521/userland/api/error/internal"
)

type Error interface {
	Error() string
	StatusCode() int
}

func NewBadRequestError(msg string) client.BadRequestError {
	return client.BadRequestError{
		Msg: msg,
	}
}

func NewUnprocessableEntityError(fields map[string]string) client.UnprocessableEntityError {
	return client.UnprocessableEntityError{
		Fields: fields,
	}
}

func NewInternalServerError() internal.InternalServerError {
	return internal.InternalServerError{}
}
