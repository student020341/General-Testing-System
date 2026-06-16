package command

import (
	"context"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/ds"
)

// calculation crud
//
//

// create

type CreateCalculationHandler struct {
	calcRepo calculation.Repository
}

func NewCreateCalculationHandler(
	cr calculation.Repository,
) CreateCalculationHandler {
	return CreateCalculationHandler{
		calcRepo: cr,
	}
}

// read

// update

type UpdateCalculationHandler struct {
	calcRepo calculation.Repository
	calcServ ds.CalculationModifiableGuard // TODO name these better
}

func NewUpdateCalculationHandler(
	cr calculation.Repository,
	cuv ds.CalculationModifiableGuard,
) UpdateCalculationHandler {
	return UpdateCalculationHandler{
		calcRepo: cr,
		calcServ: cuv,
	}
}

func (h UpdateCalculationHandler) Handle(ctx context.Context, newCalc calculation.Calculation) error {
	entity, err := h.calcRepo.GetByID(ctx, newCalc.ID)
	if err != nil {
		return err
	}

	if err := h.calcServ.EnsureCanModify(ctx, *entity); err != nil {
		return err
	}

	if err := entity.Update(calculation.UpdateCalculationFields{
		Name:    newCalc.Name,
		Closure: newCalc.Closure,
	}); err != nil {
		return err
	}

	return h.calcRepo.Save(ctx, entity)
}

// delete
