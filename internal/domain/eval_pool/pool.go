package evalpool

import (
	"errors"
	"test-system/internal/shared/optional"

	"github.com/google/uuid"
)

var (
	ErrPoolItemNotSolved = errors.New("pool item not solved")
	ErrCreateTestID      = errors.New("test id is required")
	ErrCreateEntityID    = errors.New("entity id is required")
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusComplete Status = "complete"
	StatusFailed   Status = "failed"
)

type CreatePoolItemInput struct {
	TestID          string
	EntityID        string
	DependencyCount uint
}

type PoolItem struct {
	ID              string
	TestID          string
	EntityID        string
	PoolNumber      uint
	DependencyCount uint
	Status          Status
	Err             error
}

func New(input CreatePoolItemInput) (*PoolItem, error) {
	if input.TestID == "" {
		return nil, ErrCreateTestID
	}
	if input.EntityID == "" {
		return nil, ErrCreateEntityID
	}

	return &PoolItem{
		ID:              uuid.NewString(),
		TestID:          input.TestID,
		EntityID:        input.EntityID,
		PoolNumber:      0,
		DependencyCount: input.DependencyCount,
		Status:          StatusPending,
	}, nil
}

type UpdatePoolItemInput struct {
	Status          optional.Optional[Status]
	PoolNumber      optional.Optional[uint]
	DependencyCount optional.Optional[uint]
	Err             optional.Optional[error]
}

func (p *PoolItem) Update(input UpdatePoolItemInput) {
	if input.Status.Set {
		p.Status = input.Status.Value
	}
	if input.PoolNumber.Set {
		p.PoolNumber = input.PoolNumber.Value
	}
	if input.Err.Set {
		p.Err = input.Err.Value
	}
	if input.DependencyCount.Set {
		p.DependencyCount = input.DependencyCount.Value
	}
}
