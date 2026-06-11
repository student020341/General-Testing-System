package memorymap

import (
	"errors"
	"slices"
)

var (
	ErrNotFound  = errors.New("entity not found")
	ErrInvalidID = errors.New("entity ID is invalid")
	ErrNilEntity = errors.New("cannot save nil entity")
)

type Entity interface {
	GetID() string
}

type BaseRepository[T Entity] struct {
	m   map[string]T
	iid []string
}

func NewBaseRepository[T Entity]() *BaseRepository[T] {
	return &BaseRepository[T]{
		m: make(map[string]T),
	}
}

func (r BaseRepository[T]) GetByID(
	id string,
) (T, error) {
	var zero T
	if id == "" {
		return zero, ErrInvalidID
	}
	if entry, exists := r.m[id]; exists {
		return entry, nil
	}
	return zero, ErrNotFound
}

func (r *BaseRepository[T]) Save(
	e T,
) error {
	if any(e) == nil {
		return ErrNilEntity
	}
	if e.GetID() == "" {
		return ErrInvalidID
	}

	if _, exists := r.m[e.GetID()]; !exists {
		r.iid = append(r.iid, e.GetID())
		slices.Sort(r.iid)
	}

	r.m[e.GetID()] = e
	return nil
}

func (r BaseRepository[T]) Search(
	page, pageSize int,
	matcher func(T) bool,
) (list []T, err error) {
	offset := (page - 1) * pageSize
	count := 0

	for _, id := range r.iid {
		e := r.m[id]
		if matcher(e) {
			if count < offset {
				count++
				continue
			}
			list = append(list, e)
			if len(list) >= pageSize {
				return
			}
		}
	}

	return
}

func (r *BaseRepository[T]) Delete(
	e T,
) error {
	// store concern: idempotent delete
	if any(e) == nil || e.GetID() == "" {
		return nil
	}

	r.iid = slices.DeleteFunc(r.iid, func(i string) bool { return e.GetID() == i })
	delete(r.m, e.GetID())
	return nil
}
