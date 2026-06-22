package memorymap

import (
	"context"
	"errors"
	"test-system/internal/domain/testinput"
)

var _ testinput.Repository = (*TestInputRepository)(nil)

type dbTestInput struct {
	*testinput.TestInput
}

func (t dbTestInput) GetID() string {
	return t.ID
}

func testInputFromDomain(t *testinput.TestInput) dbTestInput {
	return dbTestInput{TestInput: t}
}

func testInputToDomain(t dbTestInput) *testinput.TestInput {
	return t.TestInput
}

type TestInputRepository struct {
	*BaseRepository[dbTestInput]
}

func NewTestInputRepository() *TestInputRepository {
	return &TestInputRepository{
		BaseRepository: NewBaseRepository[dbTestInput](),
	}
}

func (r TestInputRepository) GetByID(
	ctx context.Context,
	id string,
) (*testinput.TestInput, error) {
	input, err := r.BaseRepository.GetByID(id)
	if errors.Is(err, ErrNotFound) {
		return nil, testinput.ErrNotFound
	}
	return testInputToDomain(input), err
}

func (r TestInputRepository) Search(
	ctx context.Context,
	search testinput.Search,
) ([]testinput.TestInput, error) {
	res, err := r.BaseRepository.Search(
		search.Paging,
		func(dti dbTestInput) bool {
			testMatch := search.TestID == "" || search.TestID == dti.TestID
			nameMatch := search.Name == "" || search.Name == dti.Name
			valueMatch := !search.Value.Set || search.Value.Value == dti.Value
			return testMatch && nameMatch && valueMatch
		},
	)
	if err != nil {
		return nil, err
	}

	list := make([]testinput.TestInput, len(res))
	for i, v := range res {
		d := testInputToDomain(v)
		if d == nil {
			d = &testinput.TestInput{}
		}
		list[i] = *d
	}

	return list, nil
}

func (r TestInputRepository) Save(
	ctx context.Context,
	input *testinput.TestInput,
) error {
	return r.BaseRepository.Save(testInputFromDomain(input))
}

func (r TestInputRepository) Delete(
	ctx context.Context,
	input *testinput.TestInput,
) error {
	return r.BaseRepository.Delete(testInputFromDomain(input))
}
