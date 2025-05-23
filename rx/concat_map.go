package rx

import u "alanpinder.com/rxgo/v2/utils"

func ConcatMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(downstream chan<- U, unsubscribed <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   nil, // Not used
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						sendValue: func(_ chan<- T, _ <-chan u.Never, value T) u.SelectResult {
							return drainObservable(drainObservableArgs[U]{
								source:       projection(value),
								downstream:   downstream,
								unsubscribed: unsubscribed,
							})
						},
					}
				},
			})
		})
	}
}
