package paging

import "context"

type PageFetcher[T any] func(
	ctx context.Context,
	page PageRequest,
) ([]T, error)

type Iterator[T any] struct {
	fetcher PageFetcher[T]
	page    PageRequest
	buffer  []T
	index   int
	err     error
}

func NewIterator[T any](page PageRequest, fetcher PageFetcher[T]) *Iterator[T] {
	return &Iterator[T]{
		page:    page,
		fetcher: fetcher,
		buffer:  make([]T, 0),
	}
}

func (it *Iterator[T]) Next(ctx context.Context) bool {
	if it.err != nil {
		return false
	}

	if it.index >= len(it.buffer) {
		items, err := it.fetcher(ctx, it.page)
		if err != nil {
			it.err = err
			return false
		}

		if len(items) == 0 {
			return false
		}

		it.buffer = items
		it.index = 0
		it.page.Page++
	}

	return true
}

func (it *Iterator[T]) Value() T {
	val := it.buffer[it.index]
	it.index++
	return val
}

func (it *Iterator[T]) Error() error {
	return it.err
}
