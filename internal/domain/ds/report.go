package ds

import (
	"context"
	"fmt"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
)

type ReportAddTestValidation struct {
	testRepo labtest.Repository
}

func NewReportAddTestValidation(
	tr labtest.Repository,
) ReportAddTestValidation {
	return ReportAddTestValidation{
		testRepo: tr,
	}
}

func (v ReportAddTestValidation) Validate(
	ctx context.Context,
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

	// don't assign if name conflicts
	search := labtest.Search{
		Name:     t.Name,
		PageSize: 1,
	}
	if list, err := v.testRepo.Search(ctx, search); err != nil && len(list) > 0 {
		return fmt.Errorf("a test named %q is already assigned to report %s", t.Name, r.ID)
	}

	return nil
}
