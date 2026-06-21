package report

import (
	"encoding/json"
	"fmt"
	"net/http"
	appReport "test-system/internal/app/report"
	"test-system/internal/domain/report"
)

type CreateHandler struct {
	create appReport.CreateHandler
}

func NewCreateHandler(create appReport.CreateHandler) CreateHandler {
	return CreateHandler{
		create: create,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input report.CreateReportInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	re, err := h.create.Handle(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(re); err != nil {
		fmt.Printf("create report response: %v\n", err)
		return
	}
}
