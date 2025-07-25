package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type NewUnicastObserver[T any] = func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error

type UnicastObservable[T any] struct {
	onNewObserver NewUnicastObserver[T]
}

var _ Observable[any] = (*UnicastObservable[any])(nil)

func NewUnicastObservable[T any](onNewObserver NewUnicastObserver[T]) *UnicastObservable[T] {
	return &UnicastObservable[T]{
		onNewObserver: onNewObserver,
	}
}

func (x *UnicastObservable[T]) Subscribe(ctx *Context) (<-chan T, func(), error) {

	var downstream chan T
	var sendDownstreamEndOfStream func()

	if err := u.Wrap2(NewChannel[T](ctx, 0))(&downstream, &sendDownstreamEndOfStream); err != nil {
		return nil, nil, err
	}

	var downstreamUnsubscribed chan u.Never
	var triggerDownstreamUnsubscribed func()

	if err := u.Wrap2(NewChannel[u.Never](ctx, 0))(&downstreamUnsubscribed, &triggerDownstreamUnsubscribed); err != nil {
		return nil, nil, err
	}

	GoRun(ctx, func() {

		defer triggerDownstreamUnsubscribed()
		defer sendDownstreamEndOfStream()

		err := x.onNewObserver(ctx, downstream, downstreamUnsubscribed)

		if err != nil {
			ctx.Error(err)
		}
	})

	return downstream, triggerDownstreamUnsubscribed, nil
}
