package calculation

import "errors"

var (
	ErrNilEntity   = errors.New("calculation is nil")
	ErrTestIDBlank = errors.New("test id cannot be blank")
	ErrNameBlank   = errors.New("test name cannot be blank")
)
