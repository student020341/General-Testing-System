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
