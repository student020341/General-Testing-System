package calculation

import (
	"context"
	"errors"
)

var (
	ErrInvalidSearchPage     = errors.New("search page must be >= 1")
	ErrInvalidSearchPageSize = errors.New("search page size must be between 1 and 100")
)

type Search struct {
	TestID   string
	Name     string
	Page     int
	PageSize int
}

// Validate is used to enforce business logic
func (s Search) Validate() error {
	if s.Page < 1 {
		return ErrInvalidSearchPage
	}
	if s.PageSize < 1 || s.PageSize > 100 {
		return ErrInvalidSearchPageSize
	}
	return nil
}

// WithBounds is used to enforce valid values for safety
func (s Search) WithBounds() Search {
	if s.Page < 1 {
		s.Page = 1
	}
	if s.PageSize < 1 {
		s.PageSize = 10
	} else if s.PageSize > 100 {
		s.PageSize = 100
	}
	return s
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Calculation, error)
	Search(ctx context.Context, search Search) ([]Calculation, error)
	Save(ctx context.Context, calculation *Calculation) error
	Delete(ctx context.Context, calculation *Calculation) error
}
