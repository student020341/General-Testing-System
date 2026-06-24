package labtest

import (
	"context"
	"errors"
	"fmt"
	"test-system/internal/domain/calculation"
	calculationlink "test-system/internal/domain/calculation_link"
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
}

func NewEvaluateTestHandler(
	testInputRepo testinput.Repository,
	calcRepo calculation.Repository,
	linkRepo calculationlink.Repository,
	calcLinksRepo query.CalculationsWithLinksQuery,
) EvaluateTestHandler {
	return EvaluateTestHandler{
		testInputRepo: testInputRepo,
		calcRepo:      calcRepo,
		linkRepo:      linkRepo,
		calcLinksRepo: calcLinksRepo,
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

	// phase 1: run all calculations that have no parameters
	noParamIt := paging.NewIterator(
		paging.NewPageRequest(1, 10),
		func(ctx context.Context, page paging.PageRequest) ([]calculation.Calculation, error) {
			return h.calcRepo.Search(ctx, calculation.Search{
				Paging:          page,
				TestID:          testID,
				HasDependencies: optional.New(false),
				IsSolved:        optional.New(false),
			})
		},
	)
	for noParamIt.Next(ctx) {
		calc := noParamIt.Value()
		if err := calc.Evaluate(nil); err != nil {
			return err
		}
		if err := h.calcRepo.Save(ctx, &calc); err != nil {
			return err
		}
	}
	if err := noParamIt.Error(); err != nil {
		return err
	}

	// phase 2: run all calculations that have parameters that are test inputs
	inputParamIt := paging.NewIterator(
		paging.NewPageRequest(1, 10),
		func(ctx context.Context, page paging.PageRequest) ([]query.CalculationWithLinks, error) {
			return h.calcLinksRepo.Get(
				ctx,
				page,
				testID,
				query.LinkTypeInput,
			)
		},
	)
	for inputParamIt.Next(ctx) {
		calc := inputParamIt.Value()
		args, err := calc.ToEvalInput()
		if err != nil {
			// TODO dedicated pass to add more information to all errors
			return err
		}
		if err := calc.Root.Evaluate(args); err != nil {
			return err
		}
		if err := h.calcRepo.Save(ctx, &calc.Root); err != nil {
			return err
		}
	}
	if err := inputParamIt.Error(); err != nil {
		return err
	}

	// TODO this needs to loop and avoid infinite looping
	// phase 3: run all calculations that have parameters that are other calculation results
	bothParamIt := paging.NewIterator(
		paging.NewPageRequest(1, 10),
		func(ctx context.Context, page paging.PageRequest) ([]query.CalculationWithLinks, error) {
			return h.calcLinksRepo.Get(
				ctx,
				page,
				testID,
				query.LinkTypeBoth,
			)
		},
	)
	for bothParamIt.Next(ctx) {
		calc := bothParamIt.Value()
		args, err := calc.ToEvalInput()
		if err != nil {
			return err
		}

		err = calc.Root.Evaluate(args)
		if errors.Is(err, calculation.ErrIncompleteEvalInput) {
			continue
		}
		if err != nil {
			return err
		}

		if err := h.calcRepo.Save(ctx, &calc.Root); err != nil {
			return err
		}

		// TODO restructure/optimize/etc
		// TODO outer loop that breaks if nothing changes from a pass, and bring functionality to other phases
		// something has changed - flag this page for a retry
		bothParamIt.RetryCurrentPage()
	}
	if err := bothParamIt.Error(); err != nil {
		return err
	}

	return nil
}
