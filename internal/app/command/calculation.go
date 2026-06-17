package command

import (
	"context"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/ds"

	"github.com/google/uuid"
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

func (h CreateCalculationHandler) Handle(ctx context.Context, input calculation.CreateCalculationInput) error {
	input.ID = uuid.NewString()
	entity, err := calculation.New(input)
	if err != nil {
		return err
	}

	return h.calcRepo.Save(ctx, entity)
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

	if err := entity.Update(calculation.CalculationFields{
		Name:    newCalc.Name,
		Closure: newCalc.Closure,
	}); err != nil {
		return err
	}

	return h.calcRepo.Save(ctx, entity)
}

// delete
