package calculation

import (
	"errors"
	"testing"
)

func TestClosureEvaluation(t *testing.T) {
	t.Run("should catch invalid closure", func(t *testing.T) {
		calc := Calculation{
			Closure: `(a, b) => `,
			ClosureDetails: &closureDetails{
				Parameters: []string{"a", "b"},
			},
		}

		params := EvalInput{
			"a": {Set: true, Value: 3},
			"b": {Set: true, Value: 4},
		}

		err := calc.Evaluate(params)
		if err != ErrClosureInvalid {
			t.Fatalf("unexpected error: %v", err)
		}

		calc.Closure = `1;`
		err = calc.Evaluate(params)
		if err != ErrClosureNotCallable {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("should catch wrong parameter count", func(t *testing.T) {
		calc := Calculation{
			Closure: `(a) => a`,
			ClosureDetails: &closureDetails{
				Parameters: []string{"a"},
			},
		}

		// passing 2 args when 1 expected
		params := EvalInput{
			"a": {Set: true, Value: 3},
			"b": {Set: true, Value: 4},
		}
		err := calc.Evaluate(params)
		if err != ErrParamCountMismatch {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("should handle unset parameters", func(t *testing.T) {
		calc := Calculation{
			Closure: `(a) => a`,
			ClosureDetails: &closureDetails{
				Parameters: []string{"a"},
			},
		}

		params := EvalInput{
			"a": {Set: false, Value: nil},
		}
		err := calc.Evaluate(params)
		if err == nil {
			t.Fatal("expected error for unset parameter")
		}

		if !errors.Is(err, ErrIncompleteEvalInput) {
			t.Fatalf("expected ErrIncompleteEvalInput, got: %v", err)
		}

		if calc.Result.Solved != false {
			t.Fatal("expected solved=false, got true")
		}
	})

	t.Run("should evaluate valid closure with args", func(t *testing.T) {
		calc := Calculation{
			Closure: `(a, b) => a + b`,
			ClosureDetails: &closureDetails{
				Parameters: []string{"a", "b"},
			},
		}

		params := EvalInput{
			"a": {Set: true, Value: 3},
			"b": {Set: true, Value: 4},
		}
		err := calc.Evaluate(params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if calc.Result.Solved != true {
			t.Fatal("expected solved=true, got false")
		}

		if calc.Result.Value != int64(7) {
			t.Fatalf("unexpected result: %v", calc.Result)
		}
	})
}
