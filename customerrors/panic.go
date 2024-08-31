package customerrors

import (
	"errors"
	"fmt"
)

type PanicError struct {
	Message   string
	Err       error
	Recovered any
}

func (e *PanicError) Error() string {
	return e.Message
}

func (e *PanicError) Unwrap() error {
	return e.Err
}

func NewRecoveredPanicError(r any) *PanicError {
	msg := fmt.Sprintf("panic: %v", r)
	return &PanicError{
		Message:   msg,
		Err:       errors.New(msg),
		Recovered: r,
	}
}

func NewPanicError(message string) *PanicError {
	return &PanicError{
		Message: message,
		Err:     errors.New("panic"),
	}
}
