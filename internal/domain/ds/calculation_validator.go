package ds

import (
	"context"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
)

type CalculationModifiableGuard struct {
	testRepo   labtest.Repository
	reportRepo report.Repository
}

func NewCalculationModifiableGuard(
	tr labtest.Repository,
	rr report.Repository,
) *CalculationModifiableGuard {
	return &CalculationModifiableGuard{
		testRepo:   tr,
		reportRepo: rr,
	}
}

func (v CalculationModifiableGuard) EnsureCanModify(
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
