package main

import (
	"fmt"
	"log"
	"net/http"
	appCalc "test-system/internal/app/calculation"
	appCalcLink "test-system/internal/app/calculation_link"
	appTest "test-system/internal/app/labtest"
	appReport "test-system/internal/app/report"
	appTestInput "test-system/internal/app/testinput"
	"test-system/internal/domain/ds"
	memorymap "test-system/internal/infra/store/memory_map"
	"test-system/internal/transport/http/calculation"
	calculationlink "test-system/internal/transport/http/calculation_link"
	"test-system/internal/transport/http/labtest"
	"test-system/internal/transport/http/report"
	"test-system/internal/transport/http/testinput"
)

func main() {

	calcRepo := memorymap.NewCalculationRepository()
	testRepo := memorymap.NewTestRepository()
	reportRepo := memorymap.NewReportRepository()
	calcLinkRepo := memorymap.NewCalculationLinkRepository()
	testInputRepo := memorymap.NewTestInputRepository()

	// transport wiring
	mux := http.NewServeMux()

	report.RegisterRoutes(
		mux,
		report.CommandHandlers{
			Create: appReport.NewCreateHandler(reportRepo),
		},
	)

	labtest.RegisterRoutes(
		mux,
		labtest.CommandHandlers{
			Create: appTest.NewCreateHandler(testRepo),
		},
	)

	calculation.RegisterRoutes(
		mux,
		calculation.QueryHandlers{
			Read: appCalc.NewGetByIDHandler(calcRepo),
			List: appCalc.NewListHandler(calcRepo),
		},
		calculation.CommandHandlers{
			Create: appCalc.NewCreateHandler(
				calcRepo,
				ds.NewCalculationCreate(
					calcRepo,
					testRepo,
				),
			),
			Update: appCalc.NewUpdateHandler(
				calcRepo,
				ds.NewCalculationModifiableGuard(
					testRepo,
					reportRepo,
				),
			),
		},
	)

	calculationlink.RegisterRoutes(
		mux,
		calculationlink.CommandHandlers{
			Create: appCalcLink.NewCreateHandler(
				calcLinkRepo,
				ds.NewCalculationLinkCreate(
					calcRepo,
					testRepo,
					reportRepo,
					testInputRepo,
				),
			),
		},
	)

	testinput.RegisterRoutes(
		mux,
		testinput.CommandHandlers{
			Create: appTestInput.NewCreateHandler(testInputRepo),
		},
	)

	fmt.Println("starting server on port 2000")
	if err := http.ListenAndServe(":2000", mux); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}
