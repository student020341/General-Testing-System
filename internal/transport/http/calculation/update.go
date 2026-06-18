package calculation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-system/internal/app/command"
	"test-system/internal/domain/calculation"
)

type UpdateHandler struct {
	update command.UpdateCalculationHandler
}

func NewUpdateHandler(update command.UpdateCalculationHandler) *UpdateHandler {
	return &UpdateHandler{
		update: update,
	}
}

func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input calculation.Calculation
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	res, err := h.update.Handle(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		fmt.Printf("calculation update transport error: %v\n", err)
		return
	}
}
