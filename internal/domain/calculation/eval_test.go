package calculation

import (
	"testing"
)

func TestClosureEvaluation(t *testing.T) {
	t.Run("should catch invalid closure", func(t *testing.T) {
		calc := Calculation{
			Closure: `(a, b) => `,
		}

		err := calc.Evaluate([]any{3, 4})
		if err != ErrClosureInvalid {
			t.Fatalf("unexpected error: %v", err)
		}

		calc.Closure = `1;`
		err = calc.Evaluate([]any{3, 4})
		if err != ErrClosureNotCallable {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("should catch wrong parameter count", func(t *testing.T) {
		calc := Calculation{
			Closure: `(a) => a`,
		}

		// passing 2 args when 1 expected
		err := calc.Evaluate([]any{3, 4})
		if err != ErrParamCountMismatch {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("should evaluate valid closure with args", func(t *testing.T) {
		calc := Calculation{
			Closure: `(a, b) => a + b`,
		}

		err := calc.Evaluate([]any{3, 4})
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
