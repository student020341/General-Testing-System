package report

import "errors"

var (
	ErrReportClosed = errors.New("cannot modify tests of a closed report")
	ErrNameBlank    = errors.New("report name cannot be blank")
	ErrNotFound     = errors.New("report not found")
)
