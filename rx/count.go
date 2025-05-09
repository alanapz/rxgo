package rx

import u "alanpinder.com/rxgo/v2/utils"

func Count[T any]() OperatorFunction[T, int] {
	return func(source Observable[T]) Observable[int] {
		return NewUnicastObservable(func(valuesOut chan<- int, errorsOut chan<- error, unsubscribed <-chan u.Never) {

			var count int

			drainObservable(drainObservableArgs[T]{
				source:       source,
				valuesOut:    nil, // Not used
				errorsOut:    errorsOut,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onValue: func(chan<- T, <-chan u.Never, T) u.SelectResult {
							count++
							return u.Selection(u.SelectDone(unsubscribed), u.SelectSend(valuesOut, count))
						},
					}
				},
			})
		})
	}
}
