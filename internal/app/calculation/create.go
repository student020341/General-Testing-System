package calculation

import (
	"context"
	"test-system/internal/domain/calculation"
)

type CreateCalculationService interface {
	Create(ctx context.Context, input calculation.CreateCalculationInput) (*calculation.Calculation, error)
}

type CreateHandler struct {
	calcRepo calculation.Repository
	calcServ CreateCalculationService
}

func NewCreateHandler(
	calcRepo calculation.Repository,
	calcServ CreateCalculationService,
) CreateHandler {
	return CreateHandler{
		calcRepo: calcRepo,
		calcServ: calcServ,
	}
}

func (h CreateHandler) Handle(
	ctx context.Context,
	input calculation.CreateCalculationInput,
) (*calculation.Calculation, error) {
	entity, err := h.calcServ.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	if err := h.calcRepo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
