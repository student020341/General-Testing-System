package testinput

import "errors"

var (
	ErrNotFound    = errors.New("test input not found")
	ErrTestIDBlank = errors.New("test id cannot be blank")
	ErrTypeBlank   = errors.New("type cannot be blank")
	ErrNameBlank   = errors.New("name cannot be blank")
)
