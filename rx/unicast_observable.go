package rx

import (
	"fmt"
	"runtime/debug"
)

type NewUnicastObserver[T any] = func(observer chan<- Message[T], done <-chan Void)

type UnicastObservable[T any] struct {
	onNewObserver NewUnicastObserver[T]
	stack         string
}

var _ Observable[any] = (*UnicastObservable[any])(nil)

func NewUnicastObservable[T any](onNewObserver NewUnicastObserver[T]) *UnicastObservable[T] {
	return &UnicastObservable[T]{
		onNewObserver: onNewObserver,
		stack:         string(debug.Stack()),
	}
}

func (x *UnicastObservable[T]) Subscribe() (<-chan Message[T], func()) {

	observer := NewChannel[Message[T]](fmt.Sprintf("%s: for: %s", debug.Stack(), x.stack))
	done := NewChannel[Void](fmt.Sprintf("%s: for: %s", debug.Stack(), x.stack))

	go func() {
		defer observer.Dispose()
		x.onNewObserver(observer.Channel, done.Channel)
	}()

	return observer.Channel, done.Dispose
}
