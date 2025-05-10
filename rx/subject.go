package rx

type Subject[T any] interface {
	Observable[T]
	Next(T) bool // Accepted
	EndOfStream()
	OnEndOfStream(func()) func()
}
