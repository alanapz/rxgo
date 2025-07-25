package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func Take[T any](limit uint) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error {

			var count uint

			return drainObservable(drainObservableArgs[T]{
				Context:                ctx,
				Source:                 source,
				Downstream:             downstream,
				DownstreamUnsubscribed: downstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {

					return drainObservableLoopContext[T]{

						BeforeSelection: func(*[]u.SelectItem) error {

							if count == limit {
								return ErrEndOfStream
							}

							return nil
						},

						AfterSelection: func(msg *drainObservableMessage[T]) error {

							if msg.HasValue {
								count++
							}

							return nil
						},
					}
				},
			})
		})
	}
}

func First[T any]() OperatorFunction[T, T] {
	return Take[T](1)
}
