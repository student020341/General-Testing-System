package labtest

import (
	"context"
	"test-system/internal/shared/paging"
)

type Search struct {
	ReportID string
	Name     string
	Paging   paging.PageRequest
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Test, error)
	Search(ctx context.Context, search Search) ([]Test, error)
	Save(ctx context.Context, test *Test) error
	Delete(ctx context.Context, test *Test) error
}
