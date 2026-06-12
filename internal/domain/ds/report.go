package ds

import (
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
)

type ReportAddTestValidation struct{}

func NewReportAddTestValidation() ReportAddTestValidation {
	return ReportAddTestValidation{}
}

func (v ReportAddTestValidation) Validate(
	r report.Report,
	t labtest.Test,
) error {
	// don't assign if report is not editable
	if err := r.EnsureCanModify(); err != nil {
		return err
	}

	// don't assign if test is invalid
	if err := t.EnsureValid(); err != nil {
		return err
	}

	// TODO don't assign if name conflicts

	return nil
}
