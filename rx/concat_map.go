package rx

import u "alanpinder.com/rxgo/v2/utils"

func ConcatMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(args UnicastObserverArgs[U]) {
			drainObservable(drainObservableArgs[T]{
				Environment:            args.Environment,
				Source:                 source,
				Downstream:             nil, // Not used
				DownstreamUnsubscribed: args.DownstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						sendValue: func(_ chan<- T, _ <-chan u.Never, value T) u.SelectResult {
							return drainObservable(drainObservableArgs[U]{
								Environment:            args.Environment,
								Source:                 projection(value),
								Downstream:             args.Downstream,
								DownstreamUnsubscribed: args.DownstreamUnsubscribed,
							})
						},
					}
				},
			})
		})
	}
}
