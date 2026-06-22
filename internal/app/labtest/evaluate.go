package labtest

import (
	"context"
	"fmt"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/testinput"
	"test-system/internal/shared/optional"
	"test-system/internal/shared/paging"
)

type EvaluateTestHandler struct {
	testInputRepo testinput.Repository
	calcRepo      calculation.Repository
}

func NewEvaluateTestHandler(
	testInputRepo testinput.Repository,
	calcRepo calculation.Repository,
) EvaluateTestHandler {
	return EvaluateTestHandler{
		testInputRepo: testInputRepo,
		calcRepo:      calcRepo,
	}
}

func (h EvaluateTestHandler) Handle(
	ctx context.Context,
	testID string, // TODO
) error {
	// TODO consider whether all test inputs should be set or permit partial
	// for now, enforce all are completed
	if list, err := h.testInputRepo.Search(ctx, testinput.Search{
		Paging: paging.NewPageRequest(1, 1),
		TestID: testID,
		Value: optional.Optional[any]{
			Set: true,
		},
	}); err != nil {
		return err
	} else if len(list) != 0 {
		return fmt.Errorf("test %s has incomplete inputs", testID)
	}

	// phase 1: run all calculations that have no parameters

	// phase 2: run all calculations that have parameters that are test inputs

	// phase 3: run all calculations that have parameters that are other calculation results
	// TODO this needs to loop and avoid infinite looping

	// planning: finding a calculation args
	// 1. use calculation repository to find calculations that have no inputs
	// 2. use calculation_link repository to find parameters that exist as outputs

	return nil
}
