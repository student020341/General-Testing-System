package labtest

import (
	"net/http"
	"test-system/internal/app/labtest"
)

type CommandHandlers struct {
	Create labtest.CreateHandler
}

func RegisterRoutes(mux *http.ServeMux, cmds CommandHandlers) {
	// queries

	// commands
	createHandler := NewCreateHandler(cmds.Create)

	mux.Handle("POST /tests", createHandler)
}
