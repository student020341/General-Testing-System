package calculation

import (
	"context"
	"test-system/internal/domain/calculation"

	"github.com/google/uuid"
)

type CreateHandler struct {
	calcRepo calculation.Repository
}

func NewCreateHandler(
	calcRepo calculation.Repository,
) CreateHandler {
	return CreateHandler{
		calcRepo: calcRepo,
	}
}

func (h CreateHandler) Handle(
	ctx context.Context,
	input calculation.CreateCalculationInput,
) (*calculation.Calculation, error) {
	input.ID = uuid.NewString()
	entity, err := calculation.New(input)
	if err != nil {
		return nil, err
	}

	if err := h.calcRepo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
