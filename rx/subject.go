package rx

type Subject[T any] interface {
	Observable[T]
	PostValue(T)
	PostError(error)
	PostComplete()
}
