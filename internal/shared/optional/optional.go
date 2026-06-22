package optional

type Optional[T any] struct {
	Set   bool
	Value T
}

// Zero returns an Optional with Set=true and the zero value of T
func Zero[T any]() Optional[T] {
	return Optional[T]{
		Set: true,
	}
}

func New[T any](value T) Optional[T] {
	return Optional[T]{
		Set:   true,
		Value: value,
	}
}
