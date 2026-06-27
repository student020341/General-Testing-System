package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	appCalc "test-system/internal/app/calculation"
	appCalcLink "test-system/internal/app/calculation_link"
	appTest "test-system/internal/app/labtest"
	appReport "test-system/internal/app/report"
	appTestInput "test-system/internal/app/testinput"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	evalpool "test-system/internal/domain/eval_pool"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
	"test-system/internal/domain/service"
	"test-system/internal/domain/testinput"
	memorymap "test-system/internal/infra/store/memory_map"
	calculationTransport "test-system/internal/transport/http/calculation"
	calcLinkTransport "test-system/internal/transport/http/calculation_link"
	testTransport "test-system/internal/transport/http/labtest"
	reportTransport "test-system/internal/transport/http/report"
	testInputTransport "test-system/internal/transport/http/testinput"
	"testing"
)

// test server
//
//

type TestServer struct {
	calcRepo      calculation.Repository
	testRepo      labtest.Repository
	reportRepo    report.Repository
	calcLinkRepo  calculationlink.Repository
	testInputRepo testinput.Repository
	poolsRepo     evalpool.Repository
	server        *httptest.Server
}

func setupTestServer() TestServer {
	calcRepo := memorymap.NewCalculationRepository()
	testRepo := memorymap.NewTestRepository()
	reportRepo := memorymap.NewReportRepository()
	calcLinkRepo := memorymap.NewCalculationLinkRepository()
	testInputRepo := memorymap.NewTestInputRepository()
	poolsRepo := memorymap.NewEvalPoolRepository()

	calcLinksQuery := memorymap.NewCalculationWithLinksQuery(
		*testInputRepo,
		*calcRepo,
		*calcLinkRepo,
	)

	mux := http.NewServeMux()

	reportTransport.RegisterRoutes(
		mux,
		reportTransport.CommandHandlers{
			Create: appReport.NewCreateHandler(reportRepo),
		},
	)

	testTransport.RegisterRoutes(
		mux,
		testTransport.CommandHandlers{
			Create: appTest.NewCreateHandler(testRepo),
			Build: appTest.NewBuildPoolsHandler(
				service.NewEvaluationPoolBuilder(
					calcRepo,
					testInputRepo,
					calcLinkRepo,
					poolsRepo,
					calcLinksQuery,
				),
			),
			Evaluate: appTest.NewEvaluateTestHandler(
				testInputRepo,
				calcRepo,
				calcLinkRepo,
				calcLinksQuery,
				poolsRepo,
			),
		},
	)

	calculationTransport.RegisterRoutes(
		mux,
		calculationTransport.QueryHandlers{
			Read: appCalc.NewGetByIDHandler(calcRepo),
			List: appCalc.NewListHandler(calcRepo),
		},
		calculationTransport.CommandHandlers{
			Create: appCalc.NewCreateHandler(
				calcRepo,
				service.NewCalculationCreate(
					calcRepo,
					testRepo,
				),
			),
			Update: appCalc.NewUpdateHandler(
				calcRepo,
				service.NewCalculationModifiableGuard(
					testRepo,
					reportRepo,
				),
			),
		},
	)

	calcLinkTransport.RegisterRoutes(
		mux,
		calcLinkTransport.CommandHandlers{
			Create: appCalcLink.NewCreateHandler(
				calcLinkRepo,
				service.NewCalculationLinkCreate(
					calcRepo,
					testRepo,
					reportRepo,
					testInputRepo,
				),
			),
		},
	)

	testInputTransport.RegisterRoutes(
		mux,
		testInputTransport.CommandHandlers{
			Create: appTestInput.NewCreateHandler(testInputRepo),
		},
	)

	server := httptest.NewServer(mux)

	return TestServer{
		calcRepo:      calcRepo,
		testRepo:      testRepo,
		reportRepo:    reportRepo,
		calcLinkRepo:  calcLinkRepo,
		testInputRepo: testInputRepo,
		poolsRepo:     poolsRepo,
		server:        server,
	}
}

func (ts *TestServer) close() {
	ts.server.Close()
}

type requestOption func(req *http.Request) error

func (ts TestServer) requestWithMethod(method string) requestOption {
	return func(req *http.Request) error {
		req.Method = method
		return nil
	}
}

func (ts TestServer) requestWithPath(path string) requestOption {
	return func(req *http.Request) error {
		req.URL.Path = path
		return nil
	}
}

func (ts TestServer) requestWithPayload(payload any) requestOption {
	return func(req *http.Request) error {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewReader(b))
		return nil
	}
}

func (ts TestServer) makeRequest(
	opts ...requestOption,
) (*http.Request, error) {
	req, err := http.NewRequest("GET", ts.server.URL, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type ParsedResponse[T any] struct {
	Status int
	Body   []byte
	Data   T
}

func doRequest[T any](
	client *http.Client,
	req *http.Request,
) (pr *ParsedResponse[T], err error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data T
	if len(body) > 0 {
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("parsing response %q: %w", string(body), err)
		}
	}

	return &ParsedResponse[T]{
		Status: res.StatusCode,
		Body:   body,
		Data:   data,
	}, nil
}

// test assert
//
//

// TF is TestFun
type TF struct {
	testing.TB
}

func (tf TF) Ok(err error, msg string) {
	tf.Helper()
	if err != nil {
		tf.Fatalf("%s: %v", msg, err)
	}
}

func (tf TF) Equal(expected, actual any, msg string) {
	tf.Helper()
	if expected != actual {
		tf.Fatalf("%s: expected %v, got %v", msg, expected, actual)
	}
}

func (tf TF) NotNil(actual any, msg string) {
	tf.Helper()
	if actual == nil {
		tf.Fatalf("%s: value is nil", msg)
	}
}
