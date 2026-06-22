package optional

type Optional[T any] struct {
	Set   bool
	Value T
}
