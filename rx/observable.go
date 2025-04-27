package rx

type Observable[T any] interface {
	Subscribe() (<-chan Message[T], func())
}
