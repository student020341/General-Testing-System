package query

import (
	"context"
	"test-system/internal/domain/calculation"
)

// read

type GetCalculationByIDHandler struct {
	calcRepo calculation.Repository
}

func NewGetCalculationByIDHandler(
	cr calculation.Repository,
) GetCalculationByIDHandler {
	return GetCalculationByIDHandler{
		calcRepo: cr,
	}
}

func (h GetCalculationByIDHandler) Handle(ctx context.Context, id string) (*calculation.Calculation, error) {
	return h.calcRepo.GetByID(ctx, id)
}

// list

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
