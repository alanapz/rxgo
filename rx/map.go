package rx

import u "alanpinder.com/rxgo/v2/utils"

func Map[T any, U any](projection func(T) U) OperatorFunction[T, U] {
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
							return u.Selection(u.SelectDone(args.DownstreamUnsubscribed), u.SelectSend(args.Downstream, projection(value)))
						},
					}
				},
			})
		})
	}
}

func MapTo[T any, U any](value U) OperatorFunction[T, U] {
	return Map(func(t T) U { return value })
}
