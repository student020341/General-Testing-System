package calculationlink

import "context"

type Search struct {
	Page     int
	PageSize int
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Link, error)
	Search(ctx context.Context, search Search) ([]Link, error)
	Save(ctx context.Context, link *Link) error
	Delete(ctx context.Context, link *Link) error
}
