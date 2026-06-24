package evalpool

import (
	"context"
	"test-system/internal/shared/paging"
)

type Search struct {
	TestID string
	Status Status
	Paging paging.PageRequest
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*PoolItem, error)
	Search(ctx context.Context, search Search) ([]PoolItem, error)
	Save(ctx context.Context, pool *PoolItem) error
	Delete(ctx context.Context, pool *PoolItem) error
}
