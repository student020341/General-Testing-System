package report

import (
	"context"
	"test-system/internal/domain/report"
)

type CreateHandler struct {
	reportRepo report.Repository
}

func NewCreateHandler(reportRepo report.Repository) CreateHandler {
	return CreateHandler{
		reportRepo: reportRepo,
	}
}

func (h CreateHandler) Handle(
	ctx context.Context,
	input report.CreateReportInput,
) (*report.Report, error) {
	report, err := report.New(input)
	if err != nil {
		return nil, err
	}

	if err := h.reportRepo.Save(ctx, report); err != nil {
		return nil, err
	}

	return report, nil
}
