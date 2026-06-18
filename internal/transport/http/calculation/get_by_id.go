package calculation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-system/internal/app/query"
)

type GetByIDHandler struct {
	getByID query.GetCalculationByIDHandler
}

func NewGetByIDHandler(getByID query.GetCalculationByIDHandler) *GetByIDHandler {
	return &GetByIDHandler{
		getByID: getByID,
	}
}

func (h GetByIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, err := h.getByID.Handle(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		// TODO slog
		fmt.Printf("calculation get by id transport error: %v\n", err)
		return
	}
}
