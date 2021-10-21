package client

import (
	"fmt"
	"net/http"
)

type UnprocessableEntityError struct {
	Fields map[string]string
}

func (e UnprocessableEntityError) Error() string {
	return fmt.Sprintf("Unprocessable Entity: %v", e.Fields)
}

func (e UnprocessableEntityError) StatusCode() int {
	return http.StatusUnprocessableEntity
}
