package calculationlink

import "errors"

var (
	ErrReportIDBlank        = errors.New("report id cannot be blank")
	ErrSourceIDBlank        = errors.New("source id cannot be blank")
	ErrSourceTestIDBlank    = errors.New("source test id cannot be blank")
	ErrTargetIDBlank        = errors.New("target id cannot be blank")
	ErrTargetTestIDBlank    = errors.New("target test id cannot be blank")
	ErrTestsNotInSameReport = errors.New("tests must be in the same report")
)
