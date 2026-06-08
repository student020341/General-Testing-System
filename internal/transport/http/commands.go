package http

import (
	"encoding/json"
	"net/http"
	"test-system/internal/app/command"
	"test-system/internal/domain/calculation"
)

type CalculationCommandHttpHandler struct {
	updateHandler command.UpdateCalculationHandler
}

func NewCalculationCommandHttpHandler(
	handler command.UpdateCalculationHandler,
) *CalculationCommandHttpHandler {
	return &CalculationCommandHttpHandler{
		updateHandler: handler,
	}
}

func (h CalculationCommandHttpHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input calculation.Calculation
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.updateHandler.Handle(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// TODO framework or something to handle routing and responding easier
}
