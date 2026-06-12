package labtest

import "context"

type Search struct {
	ReportID string
	Name     string
	Page     int
	PageSize int
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Test, error)
	Search(ctx context.Context, search Search) ([]Test, error)
	Save(ctx context.Context, test *Test) error
	Delete(ctx context.Context, test *Test) error
}
