package calculation

import "errors"

var (
	ErrNilEntity   = errors.New("calculation is nil")
	ErrIDBlank     = errors.New("calculation id cannot be blank")
	ErrTestIDBlank = errors.New("test id cannot be blank")
	ErrNameBlank   = errors.New("test name cannot be blank")
)
