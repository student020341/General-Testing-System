package memorymap

import (
	"context"
	calculationlink "test-system/internal/domain/calculation_link"
)

var _ calculationlink.Repository = (*CalculationLinkRepository)(nil)

type dbCalcLink struct {
	*calculationlink.Link
}

func (c dbCalcLink) GetID() string {
	return c.ID
}

func calcLinkFromDomain(c *calculationlink.Link) dbCalcLink {
	return dbCalcLink{Link: c}
}

func calcLinkToDomain(c dbCalcLink) *calculationlink.Link {
	return c.Link
}

type CalculationLinkRepository struct {
	*BaseRepository[dbCalcLink]
}

func NewCalculationLinkRepository() *CalculationLinkRepository {
	return &CalculationLinkRepository{
		BaseRepository: NewBaseRepository[dbCalcLink](),
	}
}

func (r CalculationLinkRepository) GetByID(
	ctx context.Context,
	id string,
) (*calculationlink.Link, error) {
	link, err := r.BaseRepository.GetByID(id)
	return calcLinkToDomain(link), err
}

func (r CalculationLinkRepository) Search(
	ctx context.Context,
	search calculationlink.Search,
) ([]calculationlink.Link, error) {
	res, err := r.BaseRepository.Search(
		search.Paging,
		func(l dbCalcLink) bool {
			targetMatch := search.TargetID == "" || l.Target.ID == search.TargetID
			return targetMatch
		},
	)

	var links []calculationlink.Link
	for _, link := range res {
		links = append(links, *calcLinkToDomain(link))
	}
	return links, err
}

func (r CalculationLinkRepository) Save(
	ctx context.Context,
	link *calculationlink.Link,
) error {
	return r.BaseRepository.Save(calcLinkFromDomain(link))
}

func (r CalculationLinkRepository) Delete(
	ctx context.Context,
	link *calculationlink.Link,
) error {
	return r.BaseRepository.Delete(calcLinkFromDomain(link))
}
