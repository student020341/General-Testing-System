package calculation

import (
	"slices"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
)

var _ calculationlink.LinkSourceOutputValidator = (*LinkOutputValidator)(nil)
var _ calculationlink.LinkTargetInputValidator = (*LinkTargetInputValidator)(nil)

type LinkOutputValidator struct {
	calculation.Calculation
}

func (l LinkOutputValidator) EnsureValidSourceOutput(
	source calculationlink.Source,
) error {
	switch source.OutputType {
	case "single":
		if !l.ClosureDetails.HasSingleReturn {
			return calculationlink.ErrOutputTypeNoSingleReturn
		}
	case "array":
		if !l.ClosureDetails.HasArrayReturn {
			return calculationlink.ErrOutputTypeNoArrayReturn
		}
	case "key":
		if source.OutputName == "" {
			return calculationlink.ErrOutoutTypeKeyBlank
		}
		if !slices.Contains(l.ClosureDetails.KeyedReturnFields, source.OutputName) {
			return calculationlink.ErrOutputTypeNoKeyReturn
		}
	}
	return nil
}

type LinkTargetInputValidator struct {
	calculation.Calculation
}

func (l LinkTargetInputValidator) EnsureValidTargetInput(
	target calculationlink.Target,
) error {
	if !slices.Contains(l.ClosureDetails.Parameters, target.InputName) {
		return calculationlink.ErrTargetInputNotInParamList
	}
	return nil
}
