package ds

import (
	"context"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
)

type CalculationLinkCreate struct {
	calcRepo   calculation.Repository
	testRepo   labtest.Repository
	reportRepo report.Repository
}

func NewCalculationLinkCreate(
	calcRepo calculation.Repository,
	testRepo labtest.Repository,
	reportRepo report.Repository,
) CalculationLinkCreate {
	return CalculationLinkCreate{
		calcRepo:   calcRepo,
		testRepo:   testRepo,
		reportRepo: reportRepo,
	}
}

func (c CalculationLinkCreate) Create(
	ctx context.Context,
	input calculationlink.CreateLinkInput,
) (*calculationlink.Link, error) {
	// ensure report exists
	if input.ReportID == "" {
		return nil, calculationlink.ErrReportIDBlank
	}

	if _, err := c.reportRepo.GetByID(ctx, input.ReportID); err != nil {
		return nil, err
	}

	// ensure tests exist
	if input.Source.TestID == "" {
		return nil, calculationlink.ErrSourceTestIDBlank
	}

	if input.Target.TestID == "" {
		return nil, calculationlink.ErrTargetTestIDBlank
	}

	sourceTest, err := c.testRepo.GetByID(ctx, input.Source.TestID)
	if err != nil {
		return nil, err
	}

	targetTest, err := c.testRepo.GetByID(ctx, input.Target.TestID)
	if err != nil {
		return nil, err
	}

	// ensure tests are in the same report
	if sourceTest.ReportID != targetTest.ReportID || sourceTest.ReportID != input.ReportID {
		return nil, calculationlink.ErrTestsNotInSameReport
	}

	// ensure calculations exist
	if input.Source.ID == "" {
		return nil, calculationlink.ErrSourceIDBlank
	}

	if input.Target.ID == "" {
		return nil, calculationlink.ErrTargetIDBlank
	}

	if _, err := c.calcRepo.GetByID(ctx, input.Source.ID); err != nil {
		return nil, err
	}

	if _, err := c.calcRepo.GetByID(ctx, input.Target.ID); err != nil {
		return nil, err
	}

	return calculationlink.New(input)
}
