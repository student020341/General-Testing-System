package ds

import (
	"context"
	adapter "test-system/internal/adapter/calculation"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
	"test-system/internal/domain/testinput"
)

type CalculationLinkCreate struct {
	calcRepo      calculation.Repository
	testRepo      labtest.Repository
	reportRepo    report.Repository
	testInputRepo testinput.Repository
}

func NewCalculationLinkCreate(
	calcRepo calculation.Repository,
	testRepo labtest.Repository,
	reportRepo report.Repository,
	testInputRepo testinput.Repository,
) CalculationLinkCreate {
	return CalculationLinkCreate{
		calcRepo:      calcRepo,
		testRepo:      testRepo,
		reportRepo:    reportRepo,
		testInputRepo: testInputRepo,
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

	// because the link domain does not own these concepts, the validation is done here
	// with the help of adapters

	var sourceValidator calculationlink.LinkSourceOutputValidator
	if input.Source.OutputType == string(calculationlink.OutputTypeInput) {
		ti, err := c.testInputRepo.GetByID(ctx, input.Source.ID)
		if err != nil {
			return nil, err
		}
		sourceValidator = adapter.LinkOutputTestInputValidator{
			TestInput: *ti,
		}
	} else {
		sourceCalc, err := c.calcRepo.GetByID(ctx, input.Source.ID)
		if err != nil {
			return nil, err
		}
		sourceValidator = adapter.LinkOutputCalculationValidator{
			Calculation: *sourceCalc,
		}
	}

	if err := sourceValidator.EnsureValidSourceOutput(input.Source); err != nil {
		return nil, err
	}

	targetCalc, err := c.calcRepo.GetByID(ctx, input.Target.ID)
	if err != nil {
		return nil, err
	}

	// ensure target input exists in parameter list
	targetValidator := adapter.LinkTargetInputValidator{Calculation: *targetCalc}
	if err := targetValidator.EnsureValidTargetInput(input.Target); err != nil {
		return nil, err
	}

	return calculationlink.New(input)
}
