package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	appCalc "test-system/internal/app/calculation"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/ds"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
	memorymap "test-system/internal/infra/store/memory_map"
	calculationTransport "test-system/internal/transport/http/calculation"
	"testing"
)

// test server
//
//

type TestServer struct {
	calcRepo   calculation.Repository
	testRepo   labtest.Repository
	reportRepo report.Repository
	server     *httptest.Server
}

func setupTestServer() TestServer {
	calcRepo := memorymap.NewCalculationRepository()
	testRepo := memorymap.NewTestRepository()
	reportRepo := memorymap.NewReportRepository()

	mux := http.NewServeMux()

	calculationTransport.RegisterRoutes(
		mux,
		calculationTransport.QueryHandlers{
			Read: appCalc.NewGetByIDHandler(calcRepo),
			List: appCalc.NewListHandler(calcRepo),
		},
		calculationTransport.CommandHandlers{
			Create: appCalc.NewCreateHandler(calcRepo),
			Update: appCalc.NewUpdateHandler(
				calcRepo,
				ds.NewCalculationModifiableGuard(
					testRepo,
					reportRepo,
				),
			),
		},
	)

	server := httptest.NewServer(mux)

	return TestServer{
		calcRepo:   calcRepo,
		testRepo:   testRepo,
		reportRepo: reportRepo,
		server:     server,
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
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
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
