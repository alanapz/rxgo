package rx

type OperatorFunction[T any, U any] = func(Observable[T]) Observable[U]

func Pipe[T any, U any](source Observable[T], operator OperatorFunction[T, U]) Observable[U] {
	return operator(source)
}
