package rx

import u "alanpinder.com/rxgo/v2/utils"

func Filter[T any](filter func(T) bool) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(args UnicastObserverArgs[T]) {
			drainObservable(drainObservableArgs[T]{
				Environment:            args.Environment,
				Source:                 source,
				Downstream:             args.Downstream,
				DownstreamUnsubscribed: args.DownstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onSelection: func(msg *u.SelectReceiveMessage[T]) AfterSelectionResult {

							if msg.HasValue && !filter(msg.Value) {
								return DropMessage
							}

							return ContinueMessage
						},
					}
				},
			})
		})
	}
}
