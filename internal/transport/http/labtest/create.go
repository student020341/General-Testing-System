package labtest

import (
	"encoding/json"
	"fmt"
	"net/http"
	appTest "test-system/internal/app/labtest"
	"test-system/internal/domain/labtest"
)

type CreateHandler struct {
	create appTest.CreateHandler
}

func NewCreateHandler(create appTest.CreateHandler) *CreateHandler {
	return &CreateHandler{
		create: create,
	}
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input labtest.CreateTestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	lt, err := h.create.Handle(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(lt); err != nil {
		fmt.Printf("create labtest response: %v\n", err)
		return
	}
}
