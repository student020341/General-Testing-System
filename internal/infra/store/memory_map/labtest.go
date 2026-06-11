package memorymap

import (
	"context"
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
	return testToDomain(test), err
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
