package report

import "errors"

var (
	ErrReportClosed = errors.New("cannot modify tests of a closed report")
)
