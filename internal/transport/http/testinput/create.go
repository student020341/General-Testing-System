package testinput

import (
	"encoding/json"
	"fmt"
	"net/http"
	appTestInput "test-system/internal/app/testinput"
	"test-system/internal/domain/testinput"
)

type CreateHandler struct {
	create appTestInput.CreateHandler
}

func NewCreateHandler(create appTestInput.CreateHandler) *CreateHandler {
	return &CreateHandler{
		create: create,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input testinput.TestInputCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	ti, err := h.create.Handle(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ti); err != nil {
		fmt.Printf("create test input response: %v\n", err)
		return
	}
}
