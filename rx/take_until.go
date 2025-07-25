package rx

import (
	"errors"

	u "alanpinder.com/rxgo/v2/utils"
)

var errNotifiedSignalled = errors.New("Notifier signalled")

func TakeUntil[T any](notifiers ...Observable[Void]) OperatorFunction[T, T] {

	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error {

			var notifier <-chan Void
			var unsubscribeFromNotifier func()

			if err := u.Wrap2(Pipe(Merge(notifiers...), First[Void]()).Subscribe(ctx))(&notifier, &unsubscribeFromNotifier); err != nil {
				return err
			}

			defer unsubscribeFromNotifier()

			return drainObservable(drainObservableArgs[T]{
				Context:                ctx,
				Source:                 source,
				Downstream:             downstream,
				DownstreamUnsubscribed: downstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {

					return drainObservableLoopContext[T]{

						BeforeSelection: func(items *[]u.SelectItem) error {
							u.Append(items, u.SelectDone(notifier, u.Val(errNotifiedSignalled)))
							return nil
						},
					}
				},
			})
		})
	}
}
