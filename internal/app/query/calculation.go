package query

import (
	"context"
	"test-system/internal/domain/calculation"
)

// ListCalculationsInput copies the domain.Search in case they diverge over time
type ListCalculationsInput struct {
	TestID   string
	Name     string
	Page     int
	PageSize int
}

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
	input ListCalculationsInput,
) ([]calculation.Calculation, error) {
	// TODO validate inputs per layer
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 50 {
		input.PageSize = 25
	}

	// map to domain search
	search := calculation.Search{
		TestID:   input.TestID,
		Name:     input.Name,
		Page:     input.Page,
		PageSize: input.PageSize,
	}

	return h.calcRepo.Search(ctx, search)
}
