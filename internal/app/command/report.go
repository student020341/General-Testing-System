package command

import (
	"context"
	"test-system/internal/domain/ds"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
)

type AddTestToReportHandler struct {
	reportRepo report.Repository
	testRepo   labtest.Repository
	reportServ ds.ReportAddTestValidation
}

func NewAddTestToReportHandler(
	rr report.Repository,
	tr labtest.Repository,
	rs ds.ReportAddTestValidation,
) AddTestToReportHandler {
	return AddTestToReportHandler{
		reportRepo: rr,
		testRepo:   tr,
		reportServ: rs,
	}
}

func (h AddTestToReportHandler) Handle(
	ctx context.Context,
	reportID string,
	testID string,
) error {
	// fetch entities
	r, err := h.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return err
	}

	t, err := h.testRepo.GetByID(ctx, testID)
	if err != nil {
		return err
	}

	// validate
	if err := h.reportServ.Validate(ctx, *r, *t); err != nil {
		return err
	}

	// update
	t.AssignToReport(r.ID)

	// save
	return h.testRepo.Save(ctx, t)
}
