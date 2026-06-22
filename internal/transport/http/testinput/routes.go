package testinput

import (
	"net/http"
	"test-system/internal/app/testinput"
)

type CommandHandlers struct {
	Create testinput.CreateHandler
}

func RegisterRoutes(
	mux *http.ServeMux,
	cmds CommandHandlers,
) {
	// queries

	// commands
	createHandler := NewCreateHandler(cmds.Create)

	mux.Handle("POST /test-inputs", createHandler)
}
