package labtest

import "errors"

var (
	ErrTestNotFound  = errors.New("test not found")
	ErrTestCompleted = errors.New("cannot modify calculations of a completed test")
	ErrNameInvalid   = errors.New("test name is invalid")
)
