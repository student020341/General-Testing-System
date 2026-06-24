package service

import (
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	evalpool "test-system/internal/domain/eval_pool"
	"test-system/internal/domain/testinput"
)

type EvaluationPrecompute struct {
	calcRepo      calculation.Repository
	testinputRepo testinput.Repository
	linkRepo      calculationlink.Repository
	poolRepo      evalpool.Repository
}

func NewEvaluationPrecompute(
	calcRepo calculation.Repository,
	testinputRepo testinput.Repository,
	linkRepo calculationlink.Repository,
	poolRepo evalpool.Repository,
) EvaluationPrecompute {
	return EvaluationPrecompute{
		calcRepo:      calcRepo,
		testinputRepo: testinputRepo,
		linkRepo:      linkRepo,
		poolRepo:      poolRepo,
	}
}

func (e EvaluationPrecompute) Precompute() error {
	// similar to app/labtest/evaluate, phased processing of test entities

	// iterate across all calculations, store in eval pool 0, set dep count equal to parameter count

	// phase 1: identify calculations with no parameters or parameters that are inputs, store in eval pool 1
	// phase 2: iterate over eval pool 1, identify calculations that depend on them, decrement dep count by 1
	// phase 3: identify calculations with eval pool 0 and dep count 0, store in eval pool 2
	// loop steps 2&3 until all solved, incrementing pool number on each pass, track changes and break on no change and err

	return nil
}
