package calculationlink

import "errors"

var (
	ErrReportIDBlank             = errors.New("report id cannot be blank")
	ErrSourceIDBlank             = errors.New("source id cannot be blank")
	ErrSourceTestIDBlank         = errors.New("source test id cannot be blank")
	ErrSourceOutputTypeBlank     = errors.New("source output type cannot be blank")
	ErrTargetIDBlank             = errors.New("target id cannot be blank")
	ErrTargetTestIDBlank         = errors.New("target test id cannot be blank")
	ErrTargetInputNameBlank      = errors.New("target input name cannot be blank")
	ErrTargetInputNotInParamList = errors.New("target input not found in calculation parameter list")
	ErrTestsNotInSameReport      = errors.New("tests must be in the same report")
	ErrOutputTypeNoSingleReturn  = errors.New("output type is single but calculation has no single return")
	ErrOutputTypeNoArrayReturn   = errors.New("output type is array but calculation has no array return")
	ErrOutoutTypeKeyBlank        = errors.New("map output name cannot be blank")
	ErrOutputTypeNoKeyReturn     = errors.New("output type is key but calculation does not contain the specified key")
)
