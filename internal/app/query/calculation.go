package query

import (
	"context"
	"test-system/internal/domain/calculation"
)

type ListCalculationsHandler struct {
	calcRepo calculation.Repository
}

func NewListCalculationsHandler(
	cr calculation.Repository,
) ListCalculationsHandler {
	return ListCalculationsHandler{
		calcRepo: cr,
	}
}

func (h ListCalculationsHandler) Handle(
	ctx context.Context,
	input calculation.Search,
) ([]calculation.Calculation, error) {
	// not desired at this point to error over letting infra enforce defaults over bad input
	// if err := search.Validate(); err != nil {
	// 	return nil, err
	// }

	return h.calcRepo.Search(ctx, input)
}
