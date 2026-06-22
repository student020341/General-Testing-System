package query

import (
	"context"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	"test-system/internal/shared/paging"
)

type LinkedEntity struct {
	Link   calculationlink.Link
	Entity any // TODO consider sealed interface or something here
}

type CalculationWithLinks struct {
	Root  calculation.Calculation
	Links []LinkedEntity
}

type LinkType string

const (
	LinkTypeInput       LinkType = "input"
	LinkTypeCalculation LinkType = "calculation"
	LinkTypeBoth        LinkType = "both"
)

// CalculationsWithLinksQuery fetches a list of calculations as a CalculationWithLinks struct
// in order to include links and the referenced entity
type CalculationsWithLinksQuery interface {
	Get(
		ctx context.Context,
		paging paging.PageRequest,
		testID string,
		linkType LinkType,
	) ([]CalculationWithLinks, error)
}
