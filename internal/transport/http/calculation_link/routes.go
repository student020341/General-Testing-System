package calculationlink

import (
	"net/http"
	calculationlink "test-system/internal/app/calculation_link"
)

type CommandHandlers struct {
	Create calculationlink.CreateHandler
}

func RegisterRoutes(
	mux *http.ServeMux,
	cmds CommandHandlers,
) {
	// queries

	// commands
	createHandler := NewCreateHandler(cmds.Create)

	mux.Handle("POST /calculation-link", createHandler)
}
