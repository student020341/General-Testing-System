package evalpool

import (
	"errors"
	"fmt"
	"test-system/internal/shared/optional"

	"github.com/google/uuid"
)

var (
	ErrPoolItemNotSolved       = errors.New("pool item not solved")
	ErrCreateTestID            = errors.New("test id is required")
	ErrCreateEntityID          = errors.New("entity id is required")
	ErrCreateEntityTypeInvalid = errors.New("entity type is invalid")
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusComplete Status = "complete"
	StatusFailed   Status = "failed"
)

type EntityType string

const (
	EntityTypeCalculation EntityType = "calculation"
	EntityTypeTestInput   EntityType = "test_input"
)

func (et EntityType) Validate() error {
	if et == "" {
		return fmt.Errorf("entity type is blank: %w", ErrCreateEntityTypeInvalid)
	}

	switch et {
	case EntityTypeCalculation, EntityTypeTestInput:
		return nil
	default:
		return fmt.Errorf("unexpected entity type %q: %w", et, ErrCreateEntityTypeInvalid)
	}
}

type CreatePoolItemInput struct {
	TestID          string
	EntityID        string
	EntityType      EntityType
	DependencyCount uint
}

type PoolItem struct {
	ID              string
	TestID          string
	EntityID        string
	EntityType      EntityType
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
	if err := input.EntityType.Validate(); err != nil {
		return nil, err
	}

	return &PoolItem{
		ID:              uuid.NewString(),
		TestID:          input.TestID,
		EntityID:        input.EntityID,
		EntityType:      input.EntityType,
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
