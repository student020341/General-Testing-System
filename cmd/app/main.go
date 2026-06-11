package app

import (
	"log"
	"net/http"
	"test-system/internal/domain/report"
	memorymap "test-system/internal/infra/store/memory_map"
	transport "test-system/internal/transport/http"
)

func main() {

	calcRepo := memorymap.NewCalculationRepository()
	testRepo := memorymap.NewTestRepository()
	var reportRepo report.Repository

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

	if err := http.ListenAndServe(":2000", router); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}
