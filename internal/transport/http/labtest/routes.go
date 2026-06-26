package labtest

import (
	"net/http"
	"test-system/internal/app/labtest"
)

type CommandHandlers struct {
	Create   labtest.CreateHandler
	Build    labtest.BuildPoolsHandler
	Evaluate labtest.EvaluateTestHandler
}

func RegisterRoutes(mux *http.ServeMux, cmds CommandHandlers) {
	// queries

	// commands
	createHandler := NewCreateHandler(cmds.Create)
	buildHandler := NewBuildHandler(cmds.Build)
	evaluateHandler := NewEvaluateHandler(cmds.Evaluate)

	mux.Handle("POST /tests", createHandler)
	mux.Handle("POST /tests/{id}/build", buildHandler)
	mux.Handle("POST /tests/{id}/evaluate", evaluateHandler)
}
