package rx

import u "alanpinder.com/rxgo/v2/utils"

func Count[T any]() OperatorFunction[T, int] {
	return func(source Observable[T]) Observable[int] {
		return NewUnicastObservable(func(ctx *Context, downstream chan<- int, downstreamUnsubscribed <-chan u.Never) error {

			var count int

			return drainObservable(drainObservableArgs[T]{
				Context:                ctx,
				Source:                 source,
				Downstream:             nil, // Not used
				DownstreamUnsubscribed: downstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {

					return drainObservableLoopContext[T]{

						SendValue: func(chan<- T, <-chan u.Never, T) error {

							count++

							return u.Selection(ctx,
								u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)),
								u.SelectSend(downstream, count),
							)
						},
					}
				},
			})
		})
	}
}
