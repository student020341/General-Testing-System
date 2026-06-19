package calculation

import (
	"context"
	"test-system/internal/domain/calculation"
)

type ListHandler struct {
	calcRepo calculation.Repository
}

func NewListHandler(
	cr calculation.Repository,
) ListHandler {
	return ListHandler{
		calcRepo: cr,
	}
}

func (h ListHandler) Handle(
	ctx context.Context,
	input calculation.Search,
) ([]calculation.Calculation, error) {
	// not desired at this point to error over letting infra enforce defaults over bad input
	// if err := search.Validate(); err != nil {
	// 	return nil, err
	// }

	return h.calcRepo.Search(ctx, input)
}
