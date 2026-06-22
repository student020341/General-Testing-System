package calculation

import (
	"context"
	"errors"
	"test-system/internal/shared/optional"
	"test-system/internal/shared/paging"
)

var (
	ErrInvalidSearchPage     = errors.New("search page must be >= 1")
	ErrInvalidSearchPageSize = errors.New("search page size must be between 1 and 100")
)

type Search struct {
	TestID          string
	Name            string
	HasDependencies optional.Optional[bool]
	Paging          paging.PageRequest
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Calculation, error)
	Search(ctx context.Context, search Search) ([]Calculation, error)
	Save(ctx context.Context, calculation *Calculation) error
	Delete(ctx context.Context, calculation *Calculation) error
}
