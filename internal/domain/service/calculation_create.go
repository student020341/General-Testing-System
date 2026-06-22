package service

import (
	"context"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/labtest"
)

// CalculationCreate handles the labtest and calculation business logic of creating a valid calculation
type CalculationCreate struct {
	calcRepo calculation.Repository
	testRepo labtest.Repository
}

func NewCalculationCreate(
	calcRepo calculation.Repository,
	testRepo labtest.Repository,
) CalculationCreate {
	return CalculationCreate{
		calcRepo: calcRepo,
		testRepo: testRepo,
	}
}

func (c CalculationCreate) Create(
	ctx context.Context,
	input calculation.CreateCalculationInput,
) (*calculation.Calculation, error) {
	// ensure test exists
	if input.TestID == "" {
		return nil, calculation.ErrTestIDBlank
	}

	test, err := c.testRepo.GetByID(ctx, input.TestID)
	if err != nil {
		return nil, err
	}

	if test == nil {
		return nil, labtest.ErrTestNotFound
	}

	// create calculation
	return calculation.New(input)
}
