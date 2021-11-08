package jwt

type InvalidTokenError struct{}

func NewInvalidTokenError() InvalidTokenError {
	return InvalidTokenError{}
}

func (e InvalidTokenError) Error() string {
	return "Invalid token"
}
