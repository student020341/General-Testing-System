package memorymap

import (
	"context"
	"errors"
	"test-system/internal/domain/labtest"
)

var _ labtest.Repository = (*TestRepository)(nil)

type dbTest struct {
	*labtest.Test
}

func (t dbTest) GetID() string {
	return t.ID
}

func testFromDomain(t *labtest.Test) dbTest {
	return dbTest{Test: t}
}

func testToDomain(t dbTest) *labtest.Test {
	return t.Test
}

type TestRepository struct {
	*BaseRepository[dbTest]
}

func NewTestRepository() *TestRepository {
	return &TestRepository{
		BaseRepository: NewBaseRepository[dbTest](),
	}
}

func (r TestRepository) GetByID(
	ctx context.Context,
	id string,
) (*labtest.Test, error) {
	test, err := r.BaseRepository.GetByID(id)
	if errors.Is(err, ErrNotFound) {
		return nil, labtest.ErrTestNotFound
	}
	return testToDomain(test), err
}

func (r TestRepository) Search(
	ctx context.Context,
	search labtest.Search,
) ([]labtest.Test, error) {
	res, err := r.BaseRepository.Search(
		search.Paging,
		func(dt dbTest) bool {
			nameMatch := search.Name == "" || dt.Name == search.Name
			reportMatch := search.ReportID == "" || dt.ReportID == search.ReportID
			return nameMatch && reportMatch
		},
	)
	if err != nil {
		return nil, err
	}

	list := make([]labtest.Test, len(res))
	for i, v := range res {
		d := testToDomain(v)
		if d == nil {
			d = &labtest.Test{}
		}
		list[i] = *d
	}

	return list, nil
}

func (r TestRepository) Save(
	ctx context.Context,
	test *labtest.Test,
) error {
	return r.BaseRepository.Save(testFromDomain(test))
}

func (r TestRepository) Delete(
	ctx context.Context,
	test *labtest.Test,
) error {
	return r.BaseRepository.Delete(testFromDomain(test))
}
