package report

import (
	"net/http"
	"test-system/internal/app/report"
)

type CommandHandlers struct {
	Create report.CreateHandler
}

func RegisterRoutes(mux *http.ServeMux, cmds CommandHandlers) {
	// queries

	// commands
	createHandler := NewCreateHandler(cmds.Create)

	mux.Handle("POST /reports", createHandler)
}
