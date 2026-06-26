package calculationlink

import (
	"context"
	"test-system/internal/shared/paging"
)

type Search struct {
	SourceID string
	TargetID string
	Paging   paging.PageRequest
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Link, error)
	Search(ctx context.Context, search Search) ([]Link, error)
	Save(ctx context.Context, link *Link) error
	Delete(ctx context.Context, link *Link) error
}
