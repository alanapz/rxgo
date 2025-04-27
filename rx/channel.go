package rx

import (
	"runtime/debug"
	"sync"
)

type Channel[T any] struct {
	Channel  chan T
	label    string
	stack    string
	lock     sync.Mutex
	disposed bool
}

func NewChannel[T any](label string) *Channel[T] {
	return &Channel[T]{
		Channel: make(chan T),
		label:   label,
		stack:   string(debug.Stack()),
	}
}

func (x *Channel[T]) Dispose() {
	x.lock.Lock()
	defer x.lock.Unlock()

	if !x.disposed {
		close(x.Channel)
		x.Channel = nil
		x.disposed = true
	}
}
