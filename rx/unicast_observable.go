package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type NewUnicastObserver[T any] = func(values chan<- T, errors chan<- error, unsubscribed <-chan u.Never)

type UnicastObservable[T any] struct {
	onNewObserver NewUnicastObserver[T]
}

var _ Observable[any] = (*UnicastObservable[any])(nil)

func NewUnicastObservable[T any](onNewObserver NewUnicastObserver[T]) *UnicastObservable[T] {
	return &UnicastObservable[T]{
		onNewObserver: onNewObserver,
	}
}

func (x *UnicastObservable[T]) Subscribe() (<-chan T, <-chan error, func()) {

	valuesCleanup := u.NewCleanup("UnicastObservable values channel")
	errorsCleanup := u.NewCleanup("UnicastObservable errors channel")
	unsubscribedCleanup := u.NewCleanup("UnicastObservable unsubscribed channel")

	values := u.NewChannel[T](valuesCleanup, 0)
	errors := u.NewChannel[error](errorsCleanup, 0)
	unsubscribed := u.NewChannel[u.Never](unsubscribedCleanup, 0)

	onComplete := u.NewCondition("Waiting for subscriber to complete")

	go func() {
		defer onComplete()
		defer unsubscribedCleanup.Cleanup()
		defer valuesCleanup.Cleanup()
		defer errorsCleanup.Cleanup()
		x.onNewObserver(values, errors, unsubscribed)
	}()

	return values, errors, unsubscribedCleanup.Do
}
