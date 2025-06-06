package rx

type EndOfStreamPropagationPolicy bool

const PropogateEndOfStream EndOfStreamPropagationPolicy = true
const DoNotPropogateEndOfStream EndOfStreamPropagationPolicy = false

type Subject[T any] interface {
	Observable[T]
	Next(...T) error
	EndOfStream() error
	AddSource(Observable[T], EndOfStreamPropagationPolicy)
	OnEndOfStream(func()) func()
}
