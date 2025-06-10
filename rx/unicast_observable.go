package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type UnicastObserverArgs[T any] struct {
	Environment            *RxEnvironment
	Downstream             chan<- T
	DownstreamUnsubscribed <-chan u.Never
}

type NewUnicastObserver[T any] = func(args UnicastObserverArgs[T])

type UnicastObservable[T any] struct {
	onNewObserver NewUnicastObserver[T]
}

var _ Observable[any] = (*UnicastObservable[any])(nil)

func NewUnicastObservable[T any](onNewObserver NewUnicastObserver[T]) *UnicastObservable[T] {
	return &UnicastObservable[T]{
		onNewObserver: onNewObserver,
	}
}

func (x *UnicastObservable[T]) Subscribe(env *RxEnvironment) (<-chan T, func(), error) {

	downstream, sendDownstreamEndOfStream := NewChannel[T](env, 0)
	downstreamUnsubscribed, triggerDownstreamUnsubscribed := NewChannel[u.Never](env, 0)

	u.GoRun(func() {
		defer triggerDownstreamUnsubscribed.Emit()
		defer sendDownstreamEndOfStream.Emit()
		x.onNewObserver(UnicastObserverArgs[T]{Environment: env, Downstream: downstream, DownstreamUnsubscribed: downstreamUnsubscribed})
	})

	return downstream, triggerDownstreamUnsubscribed.Emit
}
