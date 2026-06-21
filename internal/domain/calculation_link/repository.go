package calculationlink

import "context"

// TODO search

type Repository interface {
	GetByID(ctx context.Context, id string) (*Link, error)
	Save(ctx context.Context, link *Link) error
	Delete(ctx context.Context, link *Link) error
}
