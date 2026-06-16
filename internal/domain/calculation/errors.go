package calculation

import "errors"

var (
	ErrNilEntity = errors.New("calculation is nil")
	ErrTestIDNil = errors.New("test id cannot be nil")
	ErrNameBlank = errors.New("test name cannot be blank")
)
