package integration

type ValidationError struct {
	msg string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{msg: msg}
}

func (a *ValidationError) Error() string {
	return a.msg
}
