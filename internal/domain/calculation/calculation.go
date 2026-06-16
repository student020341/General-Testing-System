package calculation

type CalculationInput struct{}

type CalculationOutput struct{}

type Calculation struct {
	ID      string
	TestID  string
	Name    string
	Closure string
	Inputs  []CalculationInput
	Outputs []CalculationOutput
}

type UpdateCalculationFields struct {
	Name    string
	Closure string
}

func (c *Calculation) Update(fields UpdateCalculationFields) error {
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
	c.computeInputsAndOutputs()

	return nil
}

func (c *Calculation) computeInputsAndOutputs() {}

func (c *Calculation) AssignToTest(testID string) error {
	if testID == "" {
		return ErrTestIDNil
	}
	c.TestID = testID
	return nil
}
