package test

import (
	"slices"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
	"testing"
)

func TestCalculationE2E(t *testing.T) {
	ts := setupTestServer()
	defer ts.close()

	tf := TF{TB: t}

	var reportID string
	var testID string
	var calculationID string
	testClosure := `(a, b) => a+b`

	// TODO plan better e2e testing
	t.Run("create a report", func(t *testing.T) {
		input := report.CreateReportInput{
			Name: "Test Report",
		}

		req, err := ts.makeRequest(
			ts.requestWithMethod("POST"),
			ts.requestWithPath("/reports"),
			ts.requestWithPayload(input),
		)
		tf.Ok(err, "create request")

		res, err := doRequest[*report.Report](
			ts.server.Client(),
			req,
		)
		tf.Ok(err, "create report")
		tf.Equal(res.Status, 201, "response status")
		tf.NotNil(res.Data, "response data")
		tf.Equal(res.Data.Name, "Test Report", "report name should match")
		tf.Equal(res.Data.Status, report.StatusOpen, "report status should match")

		reportID = res.Data.ID
	})

	t.Run("create a test", func(t *testing.T) {
		input := labtest.CreateTestInput{
			ReportID: reportID,
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

	t.Run("create second calculation and link", func(t *testing.T) {
		// create second calculation
		input := calculation.CreateCalculationInput{
			TestID: testID,
			CalculationFields: calculation.CalculationFields{
				Name:    "second",
				Closure: `(c) => 2 * c`,
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
		tf.Equal(res.Status, 201, "response status")
		tf.NotNil(res.Data, "response data")

		// link to first calculation
		linkInput := calculationlink.CreateLinkInput{
			ReportID: reportID,
			Source: calculationlink.Source{
				CalculationRef: calculationlink.CalculationRef{
					ID:     calculationID,
					TestID: testID,
				},
				OutputType: "single",
			},
			Target: calculationlink.Target{
				CalculationRef: calculationlink.CalculationRef{
					ID:     res.Data.ID,
					TestID: testID,
				},
				InputName: "c",
			},
		}

		req, err = ts.makeRequest(
			ts.requestWithMethod("POST"),
			ts.requestWithPath("/calculation-link"),
			ts.requestWithPayload(linkInput),
		)
		tf.Ok(err, "create request")

		res2, err := doRequest[*calculationlink.Link](
			ts.server.Client(),
			req,
		)
		tf.Ok(err, "create link")
		tf.Equal(res2.Status, 201, "response status")
		tf.NotNil(res2.Data, "response data")
		l := res2.Data
		tf.Equal(l.Source.CalculationRef.ID, calculationID, "source calculation id")
		tf.Equal(l.Target.CalculationRef.ID, res.Data.ID, "target calculation id")
		tf.Equal(l.Source.OutputType, "single", "source output type")
		tf.Equal(l.Target.InputName, "c", "target input name")
	})

	// TODO test calculation eval and link eval
	// TODO try to make unit tests for these first?
}
