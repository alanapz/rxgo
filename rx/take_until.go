package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func TakeUntil[T any](notifiers ...Observable[Void]) OperatorFunction[T, T] {

	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(args UnicastObserverArgs[T]) {

			notifier, unsubscribeFromNotifier := Pipe(Merge(notifiers...), First[Void]()).Subscribe(args.Environment)
			defer unsubscribeFromNotifier()

			drainObservable(drainObservableArgs[T]{
				Environment:            args.Environment,
				Source:                 source,
				Downstream:             args.Downstream,
				DownstreamUnsubscribed: args.DownstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						beforeSelection: func(items *[]u.SelectItem) u.SelectResult {
							u.Append(items, u.SelectDone(notifier))
							return u.ContinueResult
						},
					}
				},
			})
		})
	}
}
