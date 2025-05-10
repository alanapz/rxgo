package rx

import u "alanpinder.com/rxgo/v2/utils"

func Map[T any, U any](projection func(T) U) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(downstream chan<- U, unsubscribed <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   nil, // Not used
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						sendValue: func(_ chan<- T, _ <-chan u.Never, value T) u.SelectResult {
							return u.Selection(u.SelectDone(unsubscribed), u.SelectSend(downstream, projection(value)))
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
