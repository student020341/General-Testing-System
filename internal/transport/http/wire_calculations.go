package http

import (
	"test-system/internal/app/command"
	"test-system/internal/app/query"
	"test-system/internal/domain/calculation"
	"test-system/internal/domain/ds"
	"test-system/internal/domain/labtest"
	"test-system/internal/domain/report"
)

type CalculationDeps struct {
	CalcRepo   calculation.Repository
	TestRepo   labtest.Repository
	ReportRepo report.Repository
}

func WireCalculations(
	deps CalculationDeps,
) (*CalculationCommandHttpHandler, *CalculationQueryHttpHandler) {
	// app parts

	// commands
	updateHandler := command.NewUpdateCalculationHandler(
		deps.CalcRepo,
		*ds.NewCalculationModifiableGuard(
			deps.TestRepo,
			deps.ReportRepo,
		),
	)

	createHandler := command.NewCreateCalculationHandler(
		deps.CalcRepo,
	)

	// queries
	listHandler := query.NewListCalculationsHandler(
		deps.CalcRepo,
	)

	// transport parts
	cmdHttp := NewCalculationCommandHttpHandler(
		createHandler,
		updateHandler,
	)
	queryHttp := NewCalculationQueryHttpHandler(
		listHandler,
	)

	return cmdHttp, queryHttp
}
