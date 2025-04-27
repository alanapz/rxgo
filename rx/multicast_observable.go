package rx

// import (
// 	"fmt"
// 	"runtime/debug"
// )

// type NewUnicastObserver[T any] = func(observer chan<- Message[T], done <-chan Void)

// type MultiObservable[T any] struct {
// 	onNewObserver NewUnicastObserver[T]
// 	stack         string

// 	observers map[int]Subscription
// }

// var _ Observable[any] = (*UnicastObservable[any])(nil)

// func NewUnicastObservable[T any](onNewObserver NewUnicastObserver[T]) *UnicastObservable[T] {
// 	return &UnicastObservable[T]{
// 		onNewObserver: onNewObserver,
// 		stack:         string(debug.Stack()),
// 	}
// }

// func (x *UnicastObservable[T]) Subscribe() (<-chan Message[T], func()) {

// 	observer, closeObserver := NewChannel[Message[T]](fmt.Sprintf("%s: for: %s", debug.Stack(), x.stack))

// 	unsubscribed =

// 		func() {
// 			x.Lock.lock()
// 		}()

// 	return observer, closeObserver

// 	observer, closeObserver := NewChannel[Message[T]](fmt.Sprintf("%s: for: %s", debug.Stack(), x.stack))

// 	done, closeDone := NewChannel[Void](fmt.Sprintf("%s: for: %s", debug.Stack(), x.stack))

// 	go func() {
// 		defer closeObserver()
// 		x.onNewObserver(observer, done)
// 	}()

// 	return observer, func() {

// 	}
// }

// func (x *UnicastObservable[T]) Unsubscribe(subscription Subscription) (<-chan Message[T], func()) {

// 	observer, closeObserver := NewChannel[Message[T]](fmt.Sprintf("%s: for: %s", debug.Stack(), x.stack))

// 	unsubscribed =

// 		func() {
// 			x.Lock.lock()
// 		}()

// 	return observer, closeObserver

// 	observer, closeObserver := NewChannel[Message[T]](fmt.Sprintf("%s: for: %s", debug.Stack(), x.stack))

// 	done, closeDone := NewChannel[Void](fmt.Sprintf("%s: for: %s", debug.Stack(), x.stack))

// 	go func() {
// 		defer closeObserver()
// 		x.onNewObserver(observer, done)
// 	}()

// 	return observer, func() {

// 	}
// }
