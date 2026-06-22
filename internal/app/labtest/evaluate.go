package labtest

import (
	"context"
	"test-system/internal/domain/calculation"
)

type EvaluateTestHandler struct {
	calcRepo calculation.Repository
}

func NewEvaluateTestHandler(calcRepo calculation.Repository) EvaluateTestHandler {
	return EvaluateTestHandler{
		calcRepo: calcRepo,
	}
}

func (h EvaluateTestHandler) Handle(
	ctx context.Context,
	testID string, // TODO
) error {
	// TODO: verify all test inputs are filled in

	// phase 1: run all calculations that have no parameters

	// phase 2: run all calculations that have parameters that are test inputs

	// phase 3: run all calculations that have parameters that are other calculation results
	// TODO this needs to loop and avoid infinite looping

	// planning: finding a calculation args
	// 1. use calculation repository to find calculations that have no inputs
	// 2. use calculation_link repository to find parameters that exist as outputs

	return nil
}
