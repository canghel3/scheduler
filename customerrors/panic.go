package customerrors

import "errors"

type Panic struct {
	Message string
	Err     error
	R       any
}

func (e *Panic) Error() string {
	return e.Message
}

func (e *Panic) Unwrap() error {
	return e.Err
}

func NewRecoveredPanic(r any) *Panic {
	return &Panic{
		Message: "recovered panic",
		Err:     errors.New("panic"),
		R:       r,
	}
}

func NewPanic(message string) *Panic {
	return &Panic{
		Message: message,
		Err:     errors.New("panic"),
	}
}
