package report

import (
	"context"
	"test-system/internal/shared/paging"
)

type Search struct {
	Status Status
	Name   string
	Paging paging.PageRequest
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Report, error)
	Search(ctx context.Context, search Search) ([]Report, error)
	Save(ctx context.Context, report *Report) error
	Delete(ctx context.Context, report *Report) error
}
