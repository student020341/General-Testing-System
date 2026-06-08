package calculation

import "context"

type Search struct {
	TestID   string
	Name     string
	Page     int
	PageSize int
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Calculation, error)
	Search(ctx context.Context, search Search) ([]Calculation, error)
	Save(ctx context.Context, calculation *Calculation) error
}
