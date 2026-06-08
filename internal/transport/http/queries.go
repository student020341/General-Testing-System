package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-system/internal/app/query"
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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	input := query.ListCalculationsInput{
		TestID:   testFilter,
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
