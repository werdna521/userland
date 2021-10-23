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

func NewUnauthorizedError(msg string) client.UnauthorizedError {
	return client.UnauthorizedError{
		Msg: msg,
	}
}

func NewNotFoundError(msg string) client.NotFoundError {
	return client.NotFoundError{
		Msg: msg,
	}
}

func NewConflictError(msg string) client.ConflictError {
	return client.ConflictError{
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
