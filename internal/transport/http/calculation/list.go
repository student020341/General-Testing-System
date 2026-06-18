package calculation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"test-system/internal/app/query"
	"test-system/internal/domain/calculation"
)

type ListHandler struct {
	list query.ListCalculationsHandler
}

func NewListHandler(appQuery query.ListCalculationsHandler) *ListHandler {
	return &ListHandler{
		list: appQuery,
	}
}

func (h ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	list, err := h.list.Handle(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		// TODO slog
		fmt.Printf("calculation list transport error: %v\n", err)
		return
	}
}
