package calculation

import (
	"net/http"
	"test-system/internal/app/calculation"
)

type QueryHandlers struct {
	Read calculation.GetByIDHandler
	List calculation.ListHandler
}

type CommandHandlers struct {
	Create calculation.CreateHandler
	Update calculation.UpdateHandler
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
