package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-system/internal/app/query"
	"test-system/internal/domain/calculation"
)

type CalculationQueryHttpHandler struct {
	listHandler query.ListCalculationsHandler
}

func NewCalculationQueryHttpHandler(
	handler query.ListCalculationsHandler,
) *CalculationQueryHttpHandler {
	return &CalculationQueryHttpHandler{
		listHandler: handler,
	}
}

func (h CalculationQueryHttpHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	testFilter := r.URL.Query().Get("test_id")
	nameFilter := r.URL.Query().Get("name")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	input := calculation.Search{
		TestID:   testFilter,
		Name:     nameFilter,
		Page:     page,
		PageSize: pageSize,
	}

	results, err := h.listHandler.Handle(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results) // TODO ?
}
