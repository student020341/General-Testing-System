package main

import (
	"fmt"
	"log"
	"net/http"
	appCalc "test-system/internal/app/calculation"
	"test-system/internal/domain/ds"
	memorymap "test-system/internal/infra/store/memory_map"
	"test-system/internal/transport/http/calculation"
)

func main() {

	calcRepo := memorymap.NewCalculationRepository()
	testRepo := memorymap.NewTestRepository()
	reportRepo := memorymap.NewReportRepository()

	// transport wiring
	mux := http.NewServeMux()

	calculation.RegisterRoutes(
		mux,
		calculation.QueryHandlers{
			Read: appCalc.NewGetByIDHandler(calcRepo),
			List: appCalc.NewListHandler(calcRepo),
		},
		calculation.CommandHandlers{
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

	fmt.Println("starting server on port 2000")
	if err := http.ListenAndServe(":2000", mux); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}
