package ds

import (
	"context"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
)

type CalculationUpdateValidator struct {
	testRepo   labtest.Repository
	reportRepo report.Repository
}

func NewCalculationUpdateValidator(
	tr labtest.Repository,
	rr report.Repository,
) *CalculationUpdateValidator {
	return &CalculationUpdateValidator{
		testRepo:   tr,
		reportRepo: rr,
	}
}

func (v CalculationUpdateValidator) ValidateUpdate(
	ctx context.Context,
	calc calculation.Calculation,
) error {
	test, err := v.testRepo.GetByID(ctx, calc.TestID)
	if err != nil {
		return err
	}
	if err := test.EnsureCanModify(); err != nil {
		return err
	}

	report, err := v.reportRepo.GetByID(ctx, test.ReportID)
	if err != nil {
		return err
	}
	if err := report.EnsureCanModify(); err != nil {
		return err
	}

	return nil
}
