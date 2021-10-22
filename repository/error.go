package repository

type NotFoundError struct{}

func NewNotFoundError() NotFoundError {
	return NotFoundError{}
}

func (e NotFoundError) Error() string {
	return "row not found"
}
