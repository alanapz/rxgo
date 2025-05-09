package rx

type Subject[T any] interface {
	Observable[T]
	PostValue(T)
	PostValuesComplete()
	PostError(error)
	PostErrorsComplete()
}
