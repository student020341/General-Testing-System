package calculation

import "fmt"

type Calculation struct {
	ID             string
	TestID         string
	Name           string
	Closure        string
	ClosureDetails *closureDetails
}

// CalculationFields is calculation fields that can be edited
type CalculationFields struct {
	Name    string
	Closure string
}

// CreateCalculationInput is the required fields for creating a calculation
type CreateCalculationInput struct {
	ID     string
	TestID string
	CalculationFields
}

// New creates a new calculation. A test ID is required. If a user wants to
// copy a calculation to another test, they need to clone the calculation from
// the source test to the destination test.
func New(input CreateCalculationInput) (*Calculation, error) {
	if input.ID == "" {
		return nil, ErrIDBlank
	}

	if input.TestID == "" {
		return nil, ErrTestIDBlank
	}

	if input.Name == "" {
		return nil, ErrNameBlank
	}

	return &Calculation{
		ID:      input.ID,
		TestID:  input.TestID,
		Name:    input.Name,
		Closure: input.Closure,
	}, nil
}

func (c *Calculation) Update(fields CalculationFields) error {
	if c == nil {
		return ErrNilEntity
	}

	// validate update ok
	if fields.Name == "" {
		return ErrNameBlank
	}

	// update
	c.Name = fields.Name
	c.Closure = fields.Closure

	// compute inputs and outputs
	cd, err := parseAndWalk(c.Closure)
	if err != nil {
		return fmt.Errorf("parsing closure: %w", err)
	}

	c.ClosureDetails = cd

	return nil
}
