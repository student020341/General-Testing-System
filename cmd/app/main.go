package main

import (
	"fmt"
	"log"
	"net/http"
	memorymap "test-system/internal/infra/store/memory_map"
	transport "test-system/internal/transport/http"
)

func main() {

	calcRepo := memorymap.NewCalculationRepository()
	testRepo := memorymap.NewTestRepository()
	reportRepo := memorymap.NewReportRepository()

	// transport wiring
	calcCmd, calcQuery := transport.WireCalculations(transport.CalculationDeps{
		CalcRepo:   calcRepo,
		TestRepo:   testRepo,
		ReportRepo: reportRepo,
	})

	router := transport.NewRouter(transport.HttpHandlers{
		CalculationCommands: calcCmd,
		CalculationQueries:  calcQuery,
	})

	fmt.Println("starting server on port 2000")
	if err := http.ListenAndServe(":2000", router); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}
