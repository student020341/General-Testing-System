package command

import (
	"context"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/ds"
)

type UpdateCalculationHandler struct {
	calcRepo    calculation.Repository
	cuValidator ds.CalculationUpdateValidator // TODO name these better
}

func NewUpdateCalculationHandler(
	cr calculation.Repository,
	cuv ds.CalculationUpdateValidator,
) UpdateCalculationHandler {
	return UpdateCalculationHandler{
		calcRepo:    cr,
		cuValidator: cuv,
	}
}

func (h UpdateCalculationHandler) Handle(ctx context.Context, calc calculation.Calculation) error {
	if err := h.cuValidator.ValidateUpdate(ctx, calc); err != nil {
		return err
	}

	return h.calcRepo.Save(ctx, &calc)
}
