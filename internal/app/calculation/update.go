package calculation

import (
	"context"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/service"
)

type UpdateHandler struct {
	calcRepo calculation.Repository
	calcServ service.CalculationModifiableGuard // TODO name these better
}

func NewUpdateHandler(
	cr calculation.Repository,
	cuv service.CalculationModifiableGuard,
) UpdateHandler {
	return UpdateHandler{
		calcRepo: cr,
		calcServ: cuv,
	}
}

func (h UpdateHandler) Handle(
	ctx context.Context,
	newCalc calculation.Calculation,
) (*calculation.Calculation, error) {
	entity, err := h.calcRepo.GetByID(ctx, newCalc.ID)
	if err != nil {
		return nil, err
	}

	if err := h.calcServ.EnsureCanModify(ctx, *entity); err != nil {
		return nil, err
	}

	if err := entity.Update(calculation.CalculationFields{
		Name:    newCalc.Name,
		Closure: newCalc.Closure,
	}); err != nil {
		return nil, err
	}

	if err := h.calcRepo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
