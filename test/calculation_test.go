package test

import (
	"slices"
	"test-system/internal/domain/calculation"
	"testing"
)

func TestCalculationE2E(t *testing.T) {
	ts := setupTestServer()
	defer ts.close()

	var entityID string
	testClosure := `(a, b) => a+b`

	t.Run("create a calculation", func(t *testing.T) {
		input := calculation.CreateCalculationInput{
			TestID: "test-id",
			CalculationFields: calculation.CalculationFields{
				Name:    "farb",
				Closure: testClosure,
			},
		}

		req, err := ts.makeRequest(
			ts.requestWithMethod("POST"),
			ts.requestWithPath("/calculations"),
			ts.requestWithPayload(input),
		)
		if err != nil {
			t.Fatalf("create request: %v", err)
		}

		res, err := doRequest[*calculation.Calculation](
			ts.server.Client(),
			req,
		)
		if err != nil {
			t.Fatalf("create calculation: %v", err)
		}

		if res.Status != 201 {
			t.Fatalf("expected status 201, got %d", res.Status)
		}

		if res.Data == nil {
			t.Fatalf("response data is nil: %s", res.Body)
		}

		if res.Data.ClosureDetails.HasSingleReturn != true {
			t.Fatal("expected closure details to indicate single return")
		}

		if res.Data.Closure != testClosure {
			t.Fatalf("closure did not save correctly: %s", res.Data.Closure)
		}

		entityID = res.Data.ID
	})

	t.Run("get calculation by id", func(t *testing.T) {
		req, err := ts.makeRequest(
			ts.requestWithPath("/calculations/" + entityID),
		)
		if err != nil {
			t.Fatalf("create request: %v", err)
		}

		res, err := doRequest[*calculation.Calculation](
			ts.server.Client(),
			req,
		)
		if err != nil {
			t.Fatalf("get calculation: %v", err)
		}

		if res.Status != 200 {
			t.Fatalf("expected status 200, got %d", res.Status)
		}

		if res.Data == nil {
			t.Fatalf("response data is nil: %s", res.Body)
		}

		if res.Data.ID != entityID {
			t.Fatalf("expected id %s, got %s", entityID, res.Data.ID)
		}

		if slices.Compare(res.Data.ClosureDetails.Parameters, []string{"a", "b"}) != 0 {
			t.Fatalf("expected parameters [a, b], got %v", res.Data.ClosureDetails.Parameters)
		}
	})
}
