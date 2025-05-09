package rx

import u "alanpinder.com/rxgo/v2/utils"

func ConcatMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(valuesOut chan<- U, errorsOut chan<- error, unsubscribed <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:       source,
				valuesOut:    nil, // Not used
				errorsOut:    errorsOut,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onValue: func(_ chan<- T, _ <-chan u.Never, value T) u.SelectResult {
							return drainObservable(drainObservableArgs[U]{
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
