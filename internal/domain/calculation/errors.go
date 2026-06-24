package calculation

import "errors"

var (
	ErrNilEntity           = errors.New("calculation is nil")
	ErrTestIDBlank         = errors.New("test id cannot be blank")
	ErrNameBlank           = errors.New("test name cannot be blank")
	ErrClosureInvalid      = errors.New("closure is invalid")
	ErrClosureNotCallable  = errors.New("closure is not callable")
	ErrParamCountMismatch  = errors.New("parameter count mismatch")
	ErrIncompleteEvalInput = errors.New("evaluation input is incomplete")
)
