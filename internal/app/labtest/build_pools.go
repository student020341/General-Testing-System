package labtest

import (
	"context"
	"test-system/internal/domain/service"
)

type BuildPoolsHandler struct {
	serv service.EvaluationPoolBuilder
}

func NewBuildPoolsHandler(serv service.EvaluationPoolBuilder) BuildPoolsHandler {
	return BuildPoolsHandler{
		serv: serv,
	}
}

func (h BuildPoolsHandler) Handle(
	ctx context.Context,
	testID string,
) error {
	return h.serv.Build(ctx, testID)
}
