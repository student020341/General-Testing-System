package labtest

import (
	"net/http"
	"test-system/internal/app/labtest"
)

type CommandHandlers struct {
	Create   labtest.CreateHandler
	Evaluate labtest.EvaluateTestHandler
}

func RegisterRoutes(mux *http.ServeMux, cmds CommandHandlers) {
	// queries

	// commands
	createHandler := NewCreateHandler(cmds.Create)
	evaluateHandler := NewEvaluateHandler(cmds.Evaluate)

	mux.Handle("POST /tests", createHandler)
	mux.Handle("POST /tests/{id}/evaluate", evaluateHandler)
}
