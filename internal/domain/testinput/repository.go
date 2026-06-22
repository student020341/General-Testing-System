package testinput

import (
	"context"
	"test-system/internal/shared/optional"
	"test-system/internal/shared/paging"
)

type Search struct {
	TestID string
	Name   string
	Value  optional.Optional[any]
	Paging paging.PageRequest
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*TestInput, error)
	Search(ctx context.Context, search Search) ([]TestInput, error)
	Save(ctx context.Context, testInput *TestInput) error
	Delete(ctx context.Context, testInput *TestInput) error
}
