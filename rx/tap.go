package rx

import u "alanpinder.com/rxgo/v2/utils"

func Tap[T any](tapValue func(T)) OperatorFunction[T, T] {
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

							if tapValue != nil && msg.HasValue {
								tapValue(msg.Value)
							}

							return ContinueMessage
						},
					}
				},
			})
		})
	}
}
