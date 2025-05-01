package rx

type Observable[T any] interface {
	Subscribe() (<-chan T, <-chan error, func())
	Pipe(func(Observable[T]) Observable[T]) Observable[T]
}
