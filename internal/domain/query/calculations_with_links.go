package query

import (
	"context"
	"fmt"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	"test-system/internal/domain/testinput"
	"test-system/internal/shared/optional"
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

// ToEvalInput transforms linked entities into a map[parameter]value
func (c CalculationWithLinks) ToEvalInput() (calculation.EvalInput, error) {
	args := make(calculation.EvalInput)

	for _, le := range c.Links {
		switch e := le.Entity.(type) {
		case *calculation.Calculation:
			args[le.Link.Target.InputName] = optional.Optional[any]{
				Set:   e.Result.Solved,
				Value: e.Result.Value,
			}
		case *testinput.TestInput:
			args[le.Link.Target.InputName] = e.Value
		default:
			return nil, fmt.Errorf("unexpected entity type: %T", e)
		}
	}

	return args, nil
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
	GetByCalculationID(
		ctx context.Context,
		calcID string,
	) (*CalculationWithLinks, error)
}
