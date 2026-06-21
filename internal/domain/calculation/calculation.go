package calculation

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/google/uuid"
)

type CalculationResult struct {
	Value  any
	Solved bool
}

type Calculation struct {
	ID             string
	TestID         string
	Name           string
	Closure        string
	ClosureDetails *closureDetails
	Result         CalculationResult
}

// CalculationFields is calculation fields that can be edited
type CalculationFields struct {
	Name    string
	Closure string
}

// CreateCalculationInput is the required fields for creating a calculation
type CreateCalculationInput struct {
	TestID string
	CalculationFields
}

// New creates a new calculation. A test ID is required. If a user wants to
// copy a calculation to another test, they need to clone the calculation from
// the source test to the destination test.
func New(input CreateCalculationInput) (*Calculation, error) {
	if input.TestID == "" {
		return nil, ErrTestIDBlank
	}

	if input.Name == "" {
		return nil, ErrNameBlank
	}

	cd, err := parseAndWalk(input.Closure)
	if err != nil {
		return nil, fmt.Errorf("parsing closure: %w", err)
	}

	return &Calculation{
		ID:             uuid.NewString(),
		TestID:         input.TestID,
		Name:           input.Name,
		Closure:        input.Closure,
		ClosureDetails: cd,
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

// TODO maybe move this, maybe make goja into infra or service, TBD
// Evaluate executes the calculation closure with the given parameters
func (c *Calculation) Evaluate(params []any) error {
	// initialize goja vm
	vm := goja.New()

	// parse closure and get callable function
	res, err := vm.RunString(c.Closure)
	if err != nil {
		return ErrClosureInvalid
	}

	callable, ok := goja.AssertFunction(res)
	if !ok {
		return ErrClosureNotCallable
	}

	// verify parameter count
	jsObj := res.ToObject(vm)
	paramCount := jsObj.Get("length").ToInteger()
	if paramCount != int64(len(params)) {
		return ErrParamCountMismatch
	}

	// convert parameters to vm args
	args := make([]goja.Value, len(params))
	for i, param := range params {
		args[i] = vm.ToValue(param)
	}

	// execute
	output, err := callable(
		goja.Undefined(), // calling context
		args...,
	)
	if err != nil {
		return err
	}

	c.Result = CalculationResult{
		Value:  output.Export(),
		Solved: true,
	}

	return nil
}
