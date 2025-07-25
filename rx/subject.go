package rx

type EndOfStreamPropagationPolicy bool

const DoNotPropogateEndOfStream EndOfStreamPropagationPolicy = false
const PropogateEndOfStream EndOfStreamPropagationPolicy = true

type Subject[T any] interface {
	Observable[T]
	Next(...T) error
	EndOfStream() error
	AddSource(Observable[T], EndOfStreamPropagationPolicy)
	OnEndOfStream(func()) func()
}
