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
