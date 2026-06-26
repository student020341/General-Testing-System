package labtest

import (
	"net/http"
	appPool "test-system/internal/app/labtest"
)

type BuildHandler struct {
	build appPool.BuildPoolsHandler
}

func NewBuildHandler(build appPool.BuildPoolsHandler) *BuildHandler {
	return &BuildHandler{
		build: build,
	}
}

func (h BuildHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	testID := r.PathValue("id")

	if err := h.build.Handle(r.Context(), testID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

	// TODO a pool summary would be cool
}
