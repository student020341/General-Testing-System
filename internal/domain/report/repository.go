package report

import "context"

type Search struct {
	Status   Status
	Name     string
	Page     int
	PageSize int
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Report, error)
	Search(ctx context.Context, search Search) ([]Report, error)
	Save(ctx context.Context, report *Report) error
}
