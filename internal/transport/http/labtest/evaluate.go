package labtest

import (
	"net/http"
	appTest "test-system/internal/app/labtest"
)

type EvaluateHandler struct {
	evaluate appTest.EvaluateTestHandler
}

func NewEvaluateHandler(evaluate appTest.EvaluateTestHandler) *EvaluateHandler {
	return &EvaluateHandler{
		evaluate: evaluate,
	}
}

func (h EvaluateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	testID := r.PathValue("id")

	if err := h.evaluate.Handle(r.Context(), testID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

	// TODO test run input and result types
}
