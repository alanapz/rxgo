package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type NewUnicastObserver[T any] = func(downstream chan<- T, unsubscribed <-chan u.Never)

type UnicastObservable[T any] struct {
	onNewObserver NewUnicastObserver[T]
}

var _ Observable[any] = (*UnicastObservable[any])(nil)

func NewUnicastObservable[T any](onNewObserver NewUnicastObserver[T]) *UnicastObservable[T] {
	return &UnicastObservable[T]{
		onNewObserver: onNewObserver,
	}
}

func (x *UnicastObservable[T]) Subscribe() (<-chan T, func()) {

	var unsubscribedCleanup, downstreamCleanup u.Event

	unsubscribed := u.NewChannel[u.Never](&unsubscribedCleanup, 0)
	downstream := u.NewChannel[T](&downstreamCleanup, 0)

	u.GoRun(func() {
		defer unsubscribedCleanup.Emit()
		defer downstreamCleanup.Emit()
		x.onNewObserver(downstream, unsubscribed)
	})

	return downstream, unsubscribedCleanup.Emit
}
