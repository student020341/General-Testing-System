package calculation

import (
	"fmt"
	"slices"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	"test-system/internal/domain/testinput"
)

var _ calculationlink.LinkSourceOutputValidator = (*LinkOutputCalculationValidator)(nil)
var _ calculationlink.LinkSourceOutputValidator = (*LinkOutputTestInputValidator)(nil)
var _ calculationlink.LinkTargetInputValidator = (*LinkTargetInputValidator)(nil)

type LinkOutputCalculationValidator struct {
	calculation.Calculation
}

func (l LinkOutputCalculationValidator) EnsureValidSourceOutput(
	source calculationlink.Source,
) error {
	if source.OutputType == string(calculationlink.OutputTypeInput) {
		return fmt.Errorf("source should be a calculation: %w", calculationlink.ErrSourceTypeInvalid)
	}

	switch source.OutputType {
	case string(calculationlink.OutputTypeSingle):
		if !l.ClosureDetails.HasSingleReturn {
			return calculationlink.ErrOutputTypeNoSingleReturn
		}
	case string(calculationlink.OutputTypeArray):
		if !l.ClosureDetails.HasArrayReturn {
			return calculationlink.ErrOutputTypeNoArrayReturn
		}
	case string(calculationlink.OutputTypeKeyed):
		if source.OutputName == "" {
			return calculationlink.ErrOutoutTypeKeyBlank
		}
		if !slices.Contains(l.ClosureDetails.KeyedReturnFields, source.OutputName) {
			return calculationlink.ErrOutputTypeNoKeyReturn
		}
	}
	return nil
}

type LinkOutputTestInputValidator struct {
	testinput.TestInput
}

func (l LinkOutputTestInputValidator) EnsureValidSourceOutput(
	source calculationlink.Source,
) error {
	if source.OutputType != string(calculationlink.OutputTypeInput) {
		return fmt.Errorf("source should be a test input: %w", calculationlink.ErrSourceTypeInvalid)
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
