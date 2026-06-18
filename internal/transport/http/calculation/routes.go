package calculation

import (
	"net/http"
	"test-system/internal/app/command"
	"test-system/internal/app/query"
)

type QueryHandlers struct {
	Read query.GetCalculationByIDHandler
	List query.ListCalculationsHandler
}

type CommandHandlers struct {
	Create command.CreateCalculationHandler
	Update command.UpdateCalculationHandler
}

func RegisterRoutes(
	mux *http.ServeMux,
	query QueryHandlers,
	cmds CommandHandlers,
) {
	// queries
	readHandler := NewGetByIDHandler(query.Read)
	listHandler := NewListHandler(query.List)

	mux.Handle("GET /calculations/{id}", readHandler)
	mux.Handle("GET /calculations", listHandler)

	// commands
	createHandler := NewCreateHandler(cmds.Create)
	updateHandler := NewUpdateHandler(cmds.Update)

	mux.Handle("POST /calculations", createHandler)
	mux.Handle("PUT /calculations", updateHandler)
}
