package labtest

import (
	"context"
	"fmt"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
	evalpool "test-system/internal/domain/eval_pool"
	"test-system/internal/domain/query"
	"test-system/internal/domain/testinput"
	"test-system/internal/shared/optional"
	"test-system/internal/shared/paging"
)

type EvaluateTestHandler struct {
	testInputRepo testinput.Repository
	calcRepo      calculation.Repository
	linkRepo      calculationlink.Repository
	calcLinksRepo query.CalculationsWithLinksQuery
	poolRepo      evalpool.Repository
}

func NewEvaluateTestHandler(
	testInputRepo testinput.Repository,
	calcRepo calculation.Repository,
	linkRepo calculationlink.Repository,
	calcLinksRepo query.CalculationsWithLinksQuery,
	evalPoolRepo evalpool.Repository,
) EvaluateTestHandler {
	return EvaluateTestHandler{
		testInputRepo: testInputRepo,
		calcRepo:      calcRepo,
		linkRepo:      linkRepo,
		calcLinksRepo: calcLinksRepo,
		poolRepo:      evalPoolRepo,
	}
}

// TODO there will probably be parts of this that should be moved to a domain service

func (h EvaluateTestHandler) Handle(
	ctx context.Context,
	testID string, // TODO
) error {
	// TODO consider whether all test inputs should be set or permit partial
	// for now, enforce all are completed
	if list, err := h.testInputRepo.Search(ctx, testinput.Search{
		Paging: paging.NewPageRequest(1, 1),
		TestID: testID,
		Value:  optional.Zero[any](),
	}); err != nil {
		return err
	} else if len(list) != 0 {
		return fmt.Errorf("test %s has incomplete inputs", testID)
	}

	// iterate and solve by pool
	currentPool := 1
	haveNext := true
	for haveNext {
		poolIT := paging.NewIterator(
			paging.NewPageRequest(1, 10),
			func(ctx context.Context, page paging.PageRequest) ([]evalpool.PoolItem, error) {
				return h.poolRepo.Search(ctx, evalpool.Search{
					Paging:     page,
					TestID:     testID,
					PoolNumber: optional.New(uint(currentPool)),
				})
			},
		)
		for poolIT.Next(ctx) {
			pi := poolIT.Value()
			// value of test inputs should already be set at this point
			if pi.EntityType == evalpool.EntityTypeTestInput {
				continue
			}

			fullCalc, err := h.calcLinksRepo.GetByCalculationID(
				ctx,
				pi.EntityID,
			)
			if err != nil {
				return fmt.Errorf("full calculation query for entity %q: %w", pi.EntityID, err)
			}

			args, err := fullCalc.ToEvalInput()
			if err != nil {
				return err
			}

			if err := fullCalc.Root.Evaluate(args); err != nil {
				return err
			}
			if err := h.calcRepo.Save(ctx, &fullCalc.Root); err != nil {
				return err
			}
		}
		if err := poolIT.Error(); err != nil {
			return err
		}

		haveNext = poolIT.AtLeastOne()
		currentPool++
	}

	return nil
}
