package calculationlink

import (
	"context"
	calculationlink "test-system/internal/domain/calculation_link"
)

type CreateLinkService interface {
	Create(ctx context.Context, input calculationlink.CreateLinkInput) (*calculationlink.Link, error)
}

type CreateHandler struct {
	linkRepo calculationlink.Repository
	linkServ CreateLinkService
}

func NewCreateHandler(
	linkRepo calculationlink.Repository,
	linkServ CreateLinkService,
) CreateHandler {
	return CreateHandler{
		linkRepo: linkRepo,
		linkServ: linkServ,
	}
}

func (h CreateHandler) Handle(
	ctx context.Context,
	input calculationlink.CreateLinkInput,
) (*calculationlink.Link, error) {
	link, err := h.linkServ.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	if err := h.linkRepo.Save(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}
