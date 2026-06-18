package calculation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-system/internal/app/command"
	"test-system/internal/domain/calculation"
)

type CreateHandler struct {
	create command.CreateCalculationHandler
}

func NewCreateHandler(create command.CreateCalculationHandler) *CreateHandler {
	return &CreateHandler{
		create: create,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input calculation.CreateCalculationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	calc, err := h.create.Handle(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(calc); err != nil {
		// TODO slog?
		fmt.Printf("create calculation response: %v\n", err)
		return
	}
}
