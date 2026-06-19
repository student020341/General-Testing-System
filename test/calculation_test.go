package test

import (
	"slices"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/labtest"
	"testing"
)

func TestCalculationE2E(t *testing.T) {
	ts := setupTestServer()
	defer ts.close()

	tf := TF{TB: t}

	var testID string
	var calculationID string
	testClosure := `(a, b) => a+b`

	// TODO plan better e2e testing
	t.Run("create a test", func(t *testing.T) {
		input := labtest.CreateTestInput{
			ReportID: "report-id",
			TestFields: labtest.TestFields{
				Name:   "e2e-test",
				Status: labtest.StatusOpen,
			},
		}

		req, err := ts.makeRequest(
			ts.requestWithMethod("POST"),
			ts.requestWithPath("/tests"),
			ts.requestWithPayload(input),
		)
		tf.Ok(err, "create request")

		res, err := doRequest[*labtest.Test](
			ts.server.Client(),
			req,
		)
		tf.Ok(err, "create test")
		tf.Equal(res.Status, 201, "response status")
		tf.NotNil(res.Data, "response data")
		tf.Equal(res.Data.Name, "e2e-test", "test name should match")
		tf.Equal(res.Data.Status, labtest.StatusOpen, "test status should match")

		testID = res.Data.ID
	})

	t.Run("create a calculation", func(t *testing.T) {
		input := calculation.CreateCalculationInput{
			TestID: testID,
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
		tf.Ok(err, "create request")

		res, err := doRequest[*calculation.Calculation](
			ts.server.Client(),
			req,
		)
		tf.Ok(err, "create calculation")
		tf.Equal(res.Status, 201, "resposne status")
		tf.NotNil(res.Data, "response data")
		tf.Equal(res.Data.ClosureDetails.HasSingleReturn, true, "closure details should indicate single return")
		tf.Equal(res.Data.Closure, testClosure, "closure should be saved correctly")

		calculationID = res.Data.ID
	})

	t.Run("get calculation by id", func(t *testing.T) {
		req, err := ts.makeRequest(
			ts.requestWithPath("/calculations/" + calculationID),
		)
		tf.Ok(err, "create request")

		res, err := doRequest[*calculation.Calculation](
			ts.server.Client(),
			req,
		)
		tf.Ok(err, "get calculation")
		tf.Equal(res.Status, 200, "expected status 200")
		tf.NotNil(res.Data, "response data is nil")
		tf.Equal(res.Data.ID, calculationID, "fetch calculation by id")
		tf.Equal(slices.Compare(res.Data.ClosureDetails.Parameters, []string{"a", "b"}), 0, "closure parameters should match")
	})
}
