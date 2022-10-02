package validation

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func Error(m string) *ValidationError {
	return &ValidationError{m}
}
