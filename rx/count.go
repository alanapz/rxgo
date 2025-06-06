package rx

import u "alanpinder.com/rxgo/v2/utils"

func Count[T any]() OperatorFunction[T, int] {
	return func(source Observable[T]) Observable[int] {
		return NewUnicastObservable(func(args UnicastObserverArgs[int]) {

			var count int

			drainObservable(drainObservableArgs[T]{
				Environment:            args.Environment,
				Source:                 source,
				Downstream:             nil, // Not used
				DownstreamUnsubscribed: args.DownstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						sendValue: func(chan<- T, <-chan u.Never, T) u.SelectResult {
							count++
							return u.Selection(u.SelectDone(args.DownstreamUnsubscribed), u.SelectSend(args.Downstream, count))
						},
					}
				},
			})
		})
	}
}
