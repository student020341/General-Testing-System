package memorymap

import (
	"context"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	"test-system/internal/domain/query"
	"test-system/internal/shared/optional"
	"test-system/internal/shared/paging"
)

// note: in a real database, this will be done very differently.
// many decisions here will be made for controlling potential scale

var _ query.CalculationsWithLinksQuery = (*CalculationWithLinksQuery)(nil)

type CalculationWithLinksQuery struct {
	testInputRepo TestInputRepository
	calcRepo      CalculationRepository
	linkRepo      CalculationLinkRepository
}

func NewCalculationWithLinksQuery(
	testInputRepo TestInputRepository,
	calcRepo CalculationRepository,
	linkRepo CalculationLinkRepository,
) CalculationWithLinksQuery {
	return CalculationWithLinksQuery{
		testInputRepo: testInputRepo,
		calcRepo:      calcRepo,
		linkRepo:      linkRepo,
	}
}

func (q CalculationWithLinksQuery) Get(
	ctx context.Context,
	pr paging.PageRequest,
	testID string,
	linkType query.LinkType,
) (list []query.CalculationWithLinks, err error) {
	offset := (pr.Page - 1) * pr.PageSize
	count := 0

	calcIt := paging.NewIterator(
		paging.NewPageRequest(1, 10),
		func(ctx context.Context, page paging.PageRequest) ([]calculation.Calculation, error) {
			return q.calcRepo.Search(ctx, calculation.Search{
				TestID:          testID,
				HasDependencies: optional.New(true),
				IsSolved:        optional.New(false),
				Paging:          page,
			})
		},
	)

	var results []query.CalculationWithLinks

	for calcIt.Next(ctx) {
		if len(results) >= int(pr.PageSize) {
			break
		}

		calc := calcIt.Value()

		linkIT := paging.NewIterator(
			paging.NewPageRequest(1, 10),
			func(ctx context.Context, page paging.PageRequest) ([]calculationlink.Link, error) {
				return q.linkRepo.Search(ctx, calculationlink.Search{
					TargetID: calc.ID,
					Paging:   page,
				})
			},
		)

		cw := query.CalculationWithLinks{
			Root: calc,
		}

		shouldFilter := false
		for linkIT.Next(ctx) {
			link := linkIT.Value()

			if linkType == query.LinkTypeInput &&
				link.Source.OutputType != string(calculationlink.OutputTypeInput) {
				shouldFilter = true
				break
			}

			if linkType == query.LinkTypeCalculation &&
				link.Source.OutputType == string(calculationlink.OutputTypeInput) {
				shouldFilter = true
				break
			}

			var le any
			if link.Source.OutputType == string(calculationlink.OutputTypeInput) {
				le, err = q.testInputRepo.GetByID(ctx, link.Source.ID)
			} else {
				le, err = q.calcRepo.GetByID(ctx, link.Source.ID)
			}

			cw.Links = append(cw.Links, query.LinkedEntity{
				Link:   link,
				Entity: le,
			})
		}
		if err := linkIT.Error(); err != nil {
			return nil, err
		}

		// pass on calculations that don't match the query
		if shouldFilter {
			continue
		}

		if count >= int(offset) {
			results = append(results, cw)
		}
		count++
	}

	if err := calcIt.Error(); err != nil {
		return nil, err
	}

	return results, nil
}

func (q CalculationWithLinksQuery) GetByCalculationID(
	ctx context.Context,
	calcID string,
) (*query.CalculationWithLinks, error) {
	calc, err := q.calcRepo.GetByID(ctx, calcID)
	if err != nil {
		return nil, err
	}

	linkIT := paging.NewIterator(
		paging.NewPageRequest(1, 10),
		func(ctx context.Context, page paging.PageRequest) ([]calculationlink.Link, error) {
			return q.linkRepo.Search(ctx, calculationlink.Search{
				TargetID: calcID,
				Paging:   page,
			})
		},
	)

	cw := query.CalculationWithLinks{
		Root: *calc,
	}

	for linkIT.Next(ctx) {
		link := linkIT.Value()

		var le any
		if link.Source.OutputType == string(calculationlink.OutputTypeInput) {
			le, err = q.testInputRepo.GetByID(ctx, link.Source.ID)
		} else {
			le, err = q.calcRepo.GetByID(ctx, link.Source.ID)
		}

		cw.Links = append(cw.Links, query.LinkedEntity{
			Link:   link,
			Entity: le,
		})
	}
	if err := linkIT.Error(); err != nil {
		return nil, err
	}

	return &cw, nil
}
