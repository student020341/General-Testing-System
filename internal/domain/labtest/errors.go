package labtest

import "errors"

var (
	ErrTestNotFound  = errors.New("test not found")
	ErrTestCompleted = errors.New("cannot modify calculations of a completed test")
	ErrNameBlank     = errors.New("test name cannot be blank")
	ErrReportIDBlank = errors.New("test report ID cannot be blank")
	ErrNilEntity     = errors.New("lab test is nil")
)
