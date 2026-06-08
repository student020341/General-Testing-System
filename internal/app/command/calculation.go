package command

import (
	"context"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
)

type UpdateCalculationHandler struct {
	calcRepo   calculation.Repository
	testRepo   labtest.Repository
	reportRepo report.Repository
}

func NewUpdateCalculationHandler(
	cr calculation.Repository,
	tr labtest.Repository,
	rr report.Repository,
) UpdateCalculationHandler {
	return UpdateCalculationHandler{
		calcRepo:   cr,
		testRepo:   tr,
		reportRepo: rr,
	}
}

func (h UpdateCalculationHandler) Handle(ctx context.Context, input calculation.Calculation) error {
	calc, err := h.calcRepo.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}

	test, err := h.testRepo.GetByID(ctx, calc.TestID)
	if err != nil {
		return err
	}
	if err := test.EnsureCanModify(); err != nil {
		return err
	}

	report, err := h.reportRepo.GetByID(ctx, test.ReportID)
	if err != nil {
		return err
	}
	if err := report.EnsureCanModify(); err != nil {
		return err
	}

	// TODO validate calculation update

	return h.calcRepo.Save(ctx, &input)
}
