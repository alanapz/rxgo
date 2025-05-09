package rx

import u "alanpinder.com/rxgo/v2/utils"

func CatchError[T any](projection func(error) Observable[T]) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, unsubscribed <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:       source,
				valuesOut:    valuesOut,
				errorsOut:    nil, // Not used
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onError: func(_ chan<- error, _ <-chan u.Never, value error) u.SelectResult {
							return drainObservable(drainObservableArgs[T]{
								source:       projection(value),
								valuesOut:    valuesOut,
								errorsOut:    errorsOut,
								unsubscribed: unsubscribed,
							})
						},
					}
				},
			})
		})
	}
}
