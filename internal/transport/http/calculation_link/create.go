package calculationlink

import (
	"encoding/json"
	"fmt"
	"net/http"
	appCLink "test-system/internal/app/calculation_link"
	calculationlink "test-system/internal/domain/calculation_link"
)

type CreateHandler struct {
	create appCLink.CreateHandler
}

func NewCreateHandler(create appCLink.CreateHandler) *CreateHandler {
	return &CreateHandler{
		create: create,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input calculationlink.CreateLinkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	link, err := h.create.Handle(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(link); err != nil {
		// TODO slog?
		fmt.Printf("create calculation link response: %v\n", err)
		return
	}
}
