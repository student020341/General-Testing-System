package main

import (
	"fmt"
	"log"
	"net/http"
	"test-system/internal/app/command"
	"test-system/internal/app/query"
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
			Read: query.NewGetCalculationByIDHandler(calcRepo),
			List: query.NewListCalculationsHandler(calcRepo),
		},
		calculation.CommandHandlers{
			Create: command.NewCreateCalculationHandler(calcRepo),
			Update: command.NewUpdateCalculationHandler(
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
