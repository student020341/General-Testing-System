package calculation

import (
	"context"
	"test-system/internal/domain/calculation"
)

type GetByIDHandler struct {
	calcRepo calculation.Repository
}

func NewGetByIDHandler(
	cr calculation.Repository,
) GetByIDHandler {
	return GetByIDHandler{
		calcRepo: cr,
	}
}

func (h GetByIDHandler) Handle(ctx context.Context, id string) (*calculation.Calculation, error) {
	return h.calcRepo.GetByID(ctx, id)
}
