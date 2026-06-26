package memorymap

import (
	"context"
	evalpool "test-system/internal/domain/eval_pool"
)

var _ evalpool.Repository = (*EvalPoolRepository)(nil)

type dbEvalPool struct {
	*evalpool.PoolItem
}

func (p dbEvalPool) GetID() string {
	return p.ID
}

func evalPoolItemFromDomain(p *evalpool.PoolItem) dbEvalPool {
	return dbEvalPool{PoolItem: p}
}

func evalPoolItemToDomain(p dbEvalPool) *evalpool.PoolItem {
	return p.PoolItem
}

type EvalPoolRepository struct {
	*BaseRepository[dbEvalPool]
}

func NewEvalPoolRepository() *EvalPoolRepository {
	return &EvalPoolRepository{
		BaseRepository: NewBaseRepository[dbEvalPool](),
	}
}

func (r *EvalPoolRepository) GetByID(
	ctx context.Context,
	id string,
) (*evalpool.PoolItem, error) {
	item, err := r.BaseRepository.GetByID(id)
	return evalPoolItemToDomain(item), err
}

func (r *EvalPoolRepository) Search(
	ctx context.Context,
	search evalpool.Search,
) ([]evalpool.PoolItem, error) {
	res, err := r.BaseRepository.Search(
		search.Paging,
		func(p dbEvalPool) bool {
			// TODO update search with optional pkg?
			testMatch := search.TestID == "" || p.TestID == search.TestID
			statusMatch := search.Status == "" || p.Status == search.Status
			idMatch := search.EntityID == "" || p.EntityID == search.EntityID
			poolMatch := !search.PoolNumber.Set || search.PoolNumber.Value == p.PoolNumber
			return testMatch && statusMatch && idMatch && poolMatch
		},
	)
	if err != nil {
		return nil, err
	}

	// transform
	list := make([]evalpool.PoolItem, len(res))
	for i, v := range res {
		d := evalPoolItemToDomain(v)
		if d == nil {
			d = &evalpool.PoolItem{}
		}
		list[i] = *d
	}

	return list, nil
}

func (r *EvalPoolRepository) Save(
	ctx context.Context,
	item *evalpool.PoolItem,
) error {
	return r.BaseRepository.Save(evalPoolItemFromDomain(item))
}

func (r *EvalPoolRepository) Delete(
	ctx context.Context,
	item *evalpool.PoolItem,
) error {
	return r.BaseRepository.Delete(evalPoolItemFromDomain(item))
}

func (r *EvalPoolRepository) DeleteAllForTest(
	ctx context.Context,
	testID string,
) error {
	r.BaseRepository.iid = make([]string, 0)
	r.BaseRepository.m = make(map[string]dbEvalPool)
	return nil
}
