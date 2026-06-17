package http

import "net/http"

type HttpHandlers struct {
	CalculationCommands *CalculationCommandHttpHandler
	CalculationQueries  *CalculationQueryHttpHandler
}

func NewRouter(handlers HttpHandlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /calculations", handlers.CalculationQueries.List)
	mux.HandleFunc("POST /calculations", handlers.CalculationCommands.Create)
	mux.HandleFunc("PUT /calculations", handlers.CalculationCommands.Update)

	return mux
}
