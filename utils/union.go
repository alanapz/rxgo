package utils

type Union[V1 any, V2 any] struct {
	V1 V1
	V2 V2
}

func NewUnionV1[T any](value T) Optional[T] {
	return Optional[T]{HasValue: true, Value: value}
}

func (x Union[V1, V2]) GetV1() {

}
