package rx

import u "alanpinder.com/rxgo/v2/utils"

type NewUnicastObserver[T any] = func(values chan<- T, errors chan<- error, done <-chan u.Never)

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

	values, cleanupValues := NewChannel[T](0)
	errors, cleanupErrors := NewChannel[error](0)
	done, cleanupDone := NewChannel[u.Never](0)

	clearCondition := AddCondition("Waiting for subscriber to complete")

	go func() {
		defer cleanupValues()
		defer cleanupErrors()
		defer clearCondition()
		x.onNewObserver(values, errors, done)
	}()

	return values, errors, cleanupDone
}
