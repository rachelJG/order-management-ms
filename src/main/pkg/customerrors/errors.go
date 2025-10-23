package customerrors

// Error is a custom error interface that extends the standard error interface
type Error interface {
	error
	StatusCode() int
	ErrorCode() string
}
