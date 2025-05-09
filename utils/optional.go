package utils

type Optional[T any] struct {
	HasValue bool
	Value    T
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{HasValue: true, Value: value}
}

func NewEmptyOptional[T any]() Optional[T] {
	return Optional[T]{HasValue: false}
}
