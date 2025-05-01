package rx

type Subject[T any] interface {
	Observable[T]
	Value(T)
	Error(error)
	Complete()
}
