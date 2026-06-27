package service

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

type EvaluationPoolBuilder struct {
	calcRepo      calculation.Repository
	testinputRepo testinput.Repository
	linkRepo      calculationlink.Repository
	poolRepo      evalpool.Repository
	calcFullQuery query.CalculationsWithLinksQuery
}

func NewEvaluationPoolBuilder(
	calcRepo calculation.Repository,
	testinputRepo testinput.Repository,
	linkRepo calculationlink.Repository,
	poolRepo evalpool.Repository,
	calcFullQuery query.CalculationsWithLinksQuery,
) EvaluationPoolBuilder {
	return EvaluationPoolBuilder{
		calcRepo:      calcRepo,
		testinputRepo: testinputRepo,
		linkRepo:      linkRepo,
		poolRepo:      poolRepo,
		calcFullQuery: calcFullQuery,
	}
}

func (e EvaluationPoolBuilder) Build(
	ctx context.Context,
	testID string,
) error {
	// TODO there will either need to be an additional pass or the logic here will need to operate at the report level instead of test
	// since there may be dependencies defined in a different test under the same report

	// clear out any past pools for test
	if err := e.poolRepo.DeleteAllForTest(ctx, testID); err != nil {
		return err
	}

	// phase 1: build initial pools

	// add test inputs to pool 1
	testinputIt := paging.NewIterator(
		paging.NewPageRequest(1, 10),
		func(ctx context.Context, page paging.PageRequest) ([]testinput.TestInput, error) {
			return e.testinputRepo.Search(ctx, testinput.Search{
				Paging: page,
				TestID: testID,
			})
		},
	)
	for testinputIt.Next(ctx) {
		testinput := testinputIt.Value()

		pi, err := evalpool.New(evalpool.CreatePoolItemInput{
			TestID:          testID,
			EntityID:        testinput.ID,
			EntityType:      evalpool.EntityTypeTestInput,
			DependencyCount: 0,
		})
		if err != nil {
			return err
		}

		pi.Update(evalpool.UpdatePoolItemInput{
			PoolNumber: optional.New[uint](1),
		})

		if err := e.poolRepo.Save(ctx, pi); err != nil {
			return err
		}
	}
	if err := testinputIt.Error(); err != nil {
		return err
	}

	// iterate across all calculations, store in eval pool 0, set dep count equal to parameter count.
	// if calculation has no parameters, put in eval pool 1
	calcIt := paging.NewIterator(
		paging.NewPageRequest(1, 10),
		func(ctx context.Context, page paging.PageRequest) ([]calculation.Calculation, error) {
			return e.calcRepo.Search(ctx, calculation.Search{
				Paging: page,
				TestID: testID,
			})
		},
	)
	for calcIt.Next(ctx) {
		calc := calcIt.Value()

		pi, err := evalpool.New(evalpool.CreatePoolItemInput{
			TestID:          testID,
			EntityID:        calc.ID,
			EntityType:      evalpool.EntityTypeCalculation,
			DependencyCount: uint(len(calc.ClosureDetails.Parameters)),
		})
		if err != nil {
			return err
		}

		// if no parameters, put in eval pool 1
		if len(calc.ClosureDetails.Parameters) == 0 {
			pi.Update(evalpool.UpdatePoolItemInput{
				PoolNumber: optional.New[uint](1),
			})
		}

		if err := e.poolRepo.Save(ctx, pi); err != nil {
			return err
		}
	}
	if err := calcIt.Error(); err != nil {
		return err
	}

	// phase 2: "optimized" query to identify calculations with test input parameters, also for pool 1
	fullIt := paging.NewIterator(
		paging.NewPageRequest(1, 10),
		func(ctx context.Context, page paging.PageRequest) ([]query.CalculationWithLinks, error) {
			return e.calcFullQuery.Get(
				ctx,
				page,
				testID,
				query.LinkTypeInput,
			)
		},
	)
	for fullIt.Next(ctx) {
		calc := fullIt.Value()

		// update pool item to be in pool 1
		list, err := e.poolRepo.Search(ctx, evalpool.Search{
			TestID:   testID,
			EntityID: calc.Root.ID,
		})
		if err != nil {
			return err
		}
		if len(list) != 1 {
			return fmt.Errorf("eval pool search did not return result for entity %q", calc.Root.ID)
		}
		pi := list[0]
		pi.Update(evalpool.UpdatePoolItemInput{
			PoolNumber: optional.New[uint](1),
		})
		if err := e.poolRepo.Save(ctx, &pi); err != nil {
			return err
		}
	}
	if err := fullIt.Error(); err != nil {
		return err
	}

	// phase 3: iterate over eval pool 1, identify calculations that depend on them, decrement dep count by 1
	var currentPool uint = 1
	someDepSolved := true
	for someDepSolved {
		someDepSolved = false
		poolIt := paging.NewIterator(
			paging.NewPageRequest(1, 10),
			func(ctx context.Context, page paging.PageRequest) ([]evalpool.PoolItem, error) {
				return e.poolRepo.Search(ctx, evalpool.Search{
					Paging:     page,
					TestID:     testID,
					PoolNumber: optional.New(currentPool),
				})
			},
		)
		for poolIt.Next(ctx) {
			pi := poolIt.Value()

			// find calculations that depend on this pool item
			linkIt := paging.NewIterator(
				paging.NewPageRequest(1, 10),
				func(ctx context.Context, page paging.PageRequest) ([]calculationlink.Link, error) {
					return e.linkRepo.Search(ctx, calculationlink.Search{
						Paging:   page,
						SourceID: pi.EntityID,
					})
				},
			)
			for linkIt.Next(ctx) {
				link := linkIt.Value()

				// look up target pool item and decrease dependency count
				targetPoolList, err := e.poolRepo.Search(ctx, evalpool.Search{
					TestID:   testID,
					EntityID: link.Target.ID,
				})
				if err != nil {
					return err
				}
				if len(targetPoolList) != 1 {
					return fmt.Errorf("expected exactly one pool item for entity %q, got %d", link.Target.ID, len(targetPoolList))
				}
				targetPi := targetPoolList[0]
				targetPi.DependencyCount -= 1
				if targetPi.DependencyCount == 0 {
					someDepSolved = true
					targetPi.PoolNumber = currentPool + 1
				}
				targetPi.Update(evalpool.UpdatePoolItemInput{
					DependencyCount: optional.New(targetPi.DependencyCount),
					PoolNumber:      optional.New(targetPi.PoolNumber),
				})
				if err := e.poolRepo.Save(ctx, &targetPi); err != nil {
					return err
				}
			}
			if err := linkIt.Error(); err != nil {
				return err
			}
		}
		if err := poolIt.Error(); err != nil {
			return err
		}
		currentPool++
	}

	return nil
}
