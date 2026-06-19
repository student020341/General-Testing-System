package labtest

import (
	"context"
	"test-system/internal/domain/labtest"
)

// TODO
type CreateTestService interface {
	Create(ctx context.Context, input labtest.CreateTestInput) (*labtest.Test, error)
}

type CreateHandler struct {
	testRepo labtest.Repository
}

func NewCreateHandler(testRepo labtest.Repository) CreateHandler {
	return CreateHandler{
		testRepo: testRepo,
	}
}

func (h CreateHandler) Handle(
	ctx context.Context,
	input labtest.CreateTestInput,
) (*labtest.Test, error) {
	entity, err := labtest.New(input)
	if err != nil {
		return nil, err
	}

	if err := h.testRepo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
