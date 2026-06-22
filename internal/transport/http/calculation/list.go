package calculation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	appCalc "test-system/internal/app/calculation"
	"test-system/internal/domain/calculation"
	"test-system/internal/shared/paging"
)

type ListHandler struct {
	list appCalc.ListHandler
}

func NewListHandler(appQuery appCalc.ListHandler) *ListHandler {
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
		TestID: testFilter,
		Name:   nameFilter,
		Paging: paging.NewPageRequest(uint(page), uint(pageSize)),
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
