package customerrors

import "errors"

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
	return &PanicError{
		Message:   "recovered panic",
		Err:       errors.New("panic"),
		Recovered: r,
	}
}

func NewPanicError(message string) *PanicError {
	return &PanicError{
		Message: message,
		Err:     errors.New("panic"),
	}
}
