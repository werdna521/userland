package error

type Error interface {
	Error() string
	StatusCode() int
}
