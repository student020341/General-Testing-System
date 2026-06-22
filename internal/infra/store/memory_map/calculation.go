package memorymap

import (
	"context"
	"test-system/internal/domain/calculation"
)

var _ calculation.Repository = (*CalculationRepository)(nil)

type dbCalc struct {
	*calculation.Calculation
}

func (c dbCalc) GetID() string {
	return c.ID
}

func calcFromDomain(c *calculation.Calculation) dbCalc {
	return dbCalc{Calculation: c}
}

func calcToDomain(c dbCalc) *calculation.Calculation {
	return c.Calculation
}

type CalculationRepository struct {
	*BaseRepository[dbCalc]
}

func NewCalculationRepository() *CalculationRepository {
	return &CalculationRepository{
		BaseRepository: NewBaseRepository[dbCalc](),
	}
}

func (r CalculationRepository) GetByID(
	ctx context.Context,
	id string,
) (*calculation.Calculation, error) {
	calc, err := r.BaseRepository.GetByID(id)
	return calcToDomain(calc), err
}

func (r CalculationRepository) Search(
	ctx context.Context,
	search calculation.Search,
) ([]calculation.Calculation, error) {
	res, err := r.BaseRepository.Search(
		search.Paging,
		func(dc dbCalc) bool {
			nameMatch := search.Name == "" || dc.Name == search.Name
			testMatch := search.TestID == "" || dc.TestID == search.TestID
			return nameMatch && testMatch
		},
	)
	if err != nil {
		return nil, err
	}

	// transform to domain
	list := make([]calculation.Calculation, len(res))
	for i, v := range res {
		d := calcToDomain(v)
		if d == nil {
			d = &calculation.Calculation{}
		}
		list[i] = *d
	}

	return list, nil
}

func (r CalculationRepository) Save(
	ctx context.Context,
	calc *calculation.Calculation,
) error {
	return r.BaseRepository.Save(calcFromDomain(calc))
}

func (r CalculationRepository) Delete(
	ctx context.Context,
	calc *calculation.Calculation,
) error {
	return r.BaseRepository.Delete(calcFromDomain(calc))
}
