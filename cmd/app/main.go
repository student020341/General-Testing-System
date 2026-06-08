package app

import (
	"log"
	"net/http"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
	transport "test-system/internal/transport/http"
)

func main() {

	// TODO
	var calcRepo calculation.Repository
	var testRepo labtest.Repository
	var reportRepo report.Repository

	// transport wiring
	calcCmd, calcQuery := transport.WireCalculations(transport.CalculationDeps{
		CalcRepo:   calcRepo,
		TestRepo:   testRepo,
		ReportRepo: reportRepo,
	})

	// TODO
	log.Fatalf("unimplemented")

	router := transport.NewRouter(transport.HttpHandlers{
		CalculationCommands: calcCmd,
		CalculationQueries:  calcQuery,
	})

	if err := http.ListenAndServe(":2000", router); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}
