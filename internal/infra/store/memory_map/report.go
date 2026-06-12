package memorymap

import (
	"context"
	"test-system/internal/domain/report"
)

var _ report.Repository = (*ReportRepository)(nil)

type dbReport struct {
	*report.Report
}

func (r dbReport) GetID() string {
	return r.ID
}

func reportFromDomain(r *report.Report) dbReport {
	return dbReport{Report: r}
}

func reportToDomain(r dbReport) *report.Report {
	return r.Report
}

type ReportRepository struct {
	*BaseRepository[dbReport]
}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{
		BaseRepository: NewBaseRepository[dbReport](),
	}
}

func (r ReportRepository) GetByID(
	ctx context.Context,
	id string,
) (*report.Report, error) {
	report, err := r.BaseRepository.GetByID(id)
	return reportToDomain(report), err
}

func (r ReportRepository) Search(
	ctx context.Context,
	search report.Search,
) ([]report.Report, error) {
	res, err := r.BaseRepository.Search(
		search.Page,
		search.PageSize,
		func(dr dbReport) bool {
			nameMatch := search.Name == "" || dr.Name == search.Name
			statusMatch := search.Status == "" || dr.Status == search.Status
			return nameMatch && statusMatch
		},
	)
	if err != nil {
		return nil, err
	}

	list := make([]report.Report, len(res))
	for i, v := range res {
		d := reportToDomain(v)
		if d == nil {
			d = &report.Report{}
		}
		list[i] = *d
	}

	return list, nil
}

func (r ReportRepository) Save(
	ctx context.Context,
	report *report.Report,
) error {
	return r.BaseRepository.Save(reportFromDomain(report))
}

func (r ReportRepository) Delete(
	ctx context.Context,
	report *report.Report,
) error {
	return r.BaseRepository.Delete(reportFromDomain(report))
}
