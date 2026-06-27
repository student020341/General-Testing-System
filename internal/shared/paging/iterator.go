package paging

import "context"

type PageFetcher[T any] func(
	ctx context.Context,
	page PageRequest,
) ([]T, error)

type Iterator[T any] struct {
	fetcher    PageFetcher[T]
	page       PageRequest
	buffer     []T
	index      int
	atLeastOne bool
	needsRetry bool
	err        error
}

func NewIterator[T any](page PageRequest, fetcher PageFetcher[T]) *Iterator[T] {
	return &Iterator[T]{
		page:    page,
		fetcher: fetcher,
		buffer:  make([]T, 0),
		index:   -1,
	}
}

// RetryCurrentPage flags that the calculations in this page window changed.
// It will re-fetch the exact page that was just processed.
func (it *Iterator[T]) RetryCurrentPage() {
	it.needsRetry = true
}

func (it *Iterator[T]) Next(ctx context.Context) bool {
	if it.err != nil {
		return false
	}

	it.index++

	if it.index >= len(it.buffer) {
		if it.needsRetry {
			it.needsRetry = false
			if it.page.Page > 0 {
				it.page.Page--
			}
		}
		return it.loadNextPage(ctx)
	}

	return true
}

// loadNextPage handles the shared data fetching logic
func (it *Iterator[T]) loadNextPage(ctx context.Context) bool {
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
	it.atLeastOne = true
	return true
}

func (it *Iterator[T]) Value() T {
	return it.buffer[it.index]
}

func (it *Iterator[T]) Error() error {
	return it.err
}

// AtLeastOne returns true if at least one item was loaded in any iteration
func (it *Iterator[T]) AtLeastOne() bool {
	return it.atLeastOne
}
