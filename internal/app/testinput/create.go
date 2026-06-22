package testinput

import (
	"context"
	"test-system/internal/domain/testinput"
)

type CreateHandler struct {
	repo testinput.Repository
}

func NewCreateHandler(repo testinput.Repository) CreateHandler {
	return CreateHandler{
		repo: repo,
	}
}

func (h CreateHandler) Handle(
	ctx context.Context,
	input testinput.TestInputCreateInput,
) (*testinput.TestInput, error) {
	entity, err := testinput.New(input)
	if err != nil {
		return nil, err
	}

	if err := h.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
