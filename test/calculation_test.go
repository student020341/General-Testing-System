package test

import (
	"context"
	"fmt"
	"slices"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	evalpool "test-system/internal/domain/eval_pool"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
	"test-system/internal/domain/testinput"
	"test-system/internal/shared/optional"
	"testing"
)

func TestCalculationE2E(t *testing.T) {
	ts := setupTestServer()
	defer ts.close()

	tf := TF{TB: t}

	var reportID string
	var testID string
	var calculationID string
	var secondCalculationID string
	var finalCalculationID string // final in terms of dependencies on other things
	var testInputID string
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
				Name:    "first",
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
		finalCalculationID = res.Data.ID

		// link to first calculation
		linkInput := calculationlink.CreateLinkInput{
			ReportID: reportID,
			Source: calculationlink.Source{
				TestEntityRef: calculationlink.TestEntityRef{
					ID:     calculationID,
					TestID: testID,
				},
				OutputType: "single",
			},
			Target: calculationlink.Target{
				TestEntityRef: calculationlink.TestEntityRef{
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
		tf.Equal(l.Source.TestEntityRef.ID, calculationID, "source calculation id")
		tf.Equal(l.Target.TestEntityRef.ID, res.Data.ID, "target calculation id")
		tf.Equal(l.Source.OutputType, "single", "source output type")
		tf.Equal(l.Target.InputName, "c", "target input name")

		secondCalculationID = res.Data.ID
	})

	t.Run("create test input and link", func(t *testing.T) {
		// create test input
		input := testinput.TestInputCreateInput{
			TestID: testID, // TODO there might not be a validator on this
			Type:   testinput.TestInputTypeVariable,
			Name:   "Something",
			Value:  optional.New[any](3),
		}

		req, err := ts.makeRequest(
			ts.requestWithMethod("POST"),
			ts.requestWithPath("/test-inputs"),
			ts.requestWithPayload(input),
		)
		tf.Ok(err, "create request")

		res, err := doRequest[*testinput.TestInput](
			ts.server.Client(),
			req,
		)
		tf.Ok(err, "create test input")
		tf.Equal(res.Status, 201, "response status")
		tf.NotNil(res.Data, "response data")

		testInputID = res.Data.ID

		// link to calculation
		linkInput := calculationlink.CreateLinkInput{
			ReportID: reportID,
			Source: calculationlink.Source{
				TestEntityRef: calculationlink.TestEntityRef{
					ID:     testInputID,
					TestID: testID,
				},
				OutputType: "input",
			},
			Target: calculationlink.Target{
				TestEntityRef: calculationlink.TestEntityRef{
					ID:     calculationID,
					TestID: testID,
				},
				InputName: "a", // TODO ensure multiple links don't target the same parameter
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
		tf.Equal(l.Source.TestEntityRef.ID, testInputID, "source test input id")
		tf.Equal(l.Target.TestEntityRef.ID, calculationID, "target calculation id")
		tf.Equal(l.Source.OutputType, "input", "source output type")
		tf.Equal(l.Target.InputName, "a", "target input name")
	})

	t.Run("create many calculations that have no dependencies", func(t *testing.T) {
		var lastCalc *calculation.Calculation
		for i := 0; i < 12; i++ {
			input := calculation.CreateCalculationInput{
				TestID: testID,
				CalculationFields: calculation.CalculationFields{
					Name:    fmt.Sprintf("No Dep Calc #%d", i+1),
					Closure: fmt.Sprintf("() => %d", i+1),
				},
			}

			req, err := ts.makeRequest(
				ts.requestWithMethod("POST"),
				ts.requestWithPath("/calculations"),
				ts.requestWithPayload(input),
			)
			tf.Ok(err, "create calculation request")

			res, err := doRequest[*calculation.Calculation](
				ts.server.Client(),
				req,
			)
			tf.Ok(err, "create calculation")
			tf.Equal(res.Status, 201, "response status")
			lastCalc = res.Data
		}

		// set last calc as input B of first calculation
		input := calculationlink.CreateLinkInput{
			ReportID: reportID,
			Source: calculationlink.Source{
				TestEntityRef: calculationlink.TestEntityRef{
					ID:     lastCalc.ID,
					TestID: testID,
				},
				OutputType: "single",
			},
			Target: calculationlink.Target{
				TestEntityRef: calculationlink.TestEntityRef{
					ID:     calculationID,
					TestID: testID,
				},
				InputName: "b",
			},
		}

		req, err := ts.makeRequest(
			ts.requestWithMethod("POST"),
			ts.requestWithPath("/calculation-link"),
			ts.requestWithPayload(input),
		)
		tf.Ok(err, "create link request")

		res, err := doRequest[*calculationlink.Link](
			ts.server.Client(),
			req,
		)
		tf.Ok(err, "create link")
		tf.Equal(res.Status, 201, "response status")
		tf.NotNil(res.Data, "response data")
	})

	t.Run("test build evaluation pools", func(t *testing.T) {
		req, err := ts.makeRequest(
			ts.requestWithMethod("POST"),
			ts.requestWithPath("/tests/"+testID+"/build"),
		)
		tf.Ok(err, "build request")

		res, err := doRequest[any](
			ts.server.Client(),
			req,
		)
		tf.Ok(err, "build test")
		tf.Equal(res.Status, 204, "response status")

		// calc 1 eval pool should be 2
		piList, err := ts.poolsRepo.Search(context.Background(), evalpool.Search{
			EntityID: calculationID,
		})
		tf.Ok(err, "search calc 1 eval pool")
		tf.Equal(len(piList), 1, "fetch calc 1 pool item")
		pi := piList[0]
		tf.Equal(pi.PoolNumber, uint(2), "calculation 1 is in pool 2")

		// calc 2 eval pool should be 3
		piList, err = ts.poolsRepo.Search(context.Background(), evalpool.Search{
			EntityID: secondCalculationID,
		})
		tf.Ok(err, "search calc 2 eval pool")
		tf.Equal(len(piList), 1, "fetch calc 2 pool item")
		pi = piList[0]
		tf.Equal(pi.PoolNumber, uint(3), "calculation 2 is in pool 3")
	})

	t.Run("test eval", func(t *testing.T) {
		req, err := ts.makeRequest(
			ts.requestWithMethod("POST"),
			ts.requestWithPath("/tests/"+testID+"/evaluate"),
		)
		tf.Ok(err, "evaluate request")

		res, err := doRequest[any](
			ts.server.Client(),
			req,
		)
		tf.Ok(err, "evaluate test")
		tf.Equal(res.Status, 204, "response status")

		// check if the last calculation is solved
		calc, err := ts.calcRepo.GetByID(context.Background(), finalCalculationID)
		tf.Ok(err, "get last calculation")
		tf.NotNil(calc, "last calculation")
		tf.Equal(calc.Result.Solved, true, "calculation is solved")

		// calculation path was:
		// (c) => 2 * c
		// c input = output of (a, b) => a+b
		// a input = test input, set to 3
		// b input = calculation defined as () => 12
		// final eval: (c=15) = 2 * 15 = 30
		tf.Equal(calc.Result.Value, int64(30), "verify calculation result")
	})

	// t.Run("scratch", func(t *testing.T) {
	// 	// temporary space while figuring out what's next

	// 	links, err := ts.calcLinkRepo.Search(
	// 		context.Background(),
	// 		calculationlink.Search{
	// 			Page:     1,
	// 			PageSize: 10,
	// 		},
	// 	)
	// 	tf.Ok(err, "search links")

	// 	jb, _ := json.MarshalIndent(links, "", "  ")
	// 	fmt.Println("links", string(jb))
	// })

	// TODO test calculation eval and link eval
	// TODO try to make unit tests for these first?
}
