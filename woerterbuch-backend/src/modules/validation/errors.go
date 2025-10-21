package validation

import "errors"

type ValidationError struct {
	Message string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{message}
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) IsValidationError() bool {
	var target *ValidationError
	return errors.As(e, &target)
}
