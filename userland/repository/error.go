package repository

type NotFoundError struct{}

func NewNotFoundError() NotFoundError {
	return NotFoundError{}
}

func (e NotFoundError) Error() string {
	return "row not found"
}

type UniqueViolationError struct{}

func NewUniqueViolationError() UniqueViolationError {
	return UniqueViolationError{}
}

func (e UniqueViolationError) Error() string {
	return "unique constraint violation"
}
