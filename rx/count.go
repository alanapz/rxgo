package rx

import u "alanpinder.com/rxgo/v2/utils"

func Count[T any]() OperatorFunction[T, int] {
	return func(source Observable[T]) Observable[int] {
		return NewUnicastObservable(func(downstream chan<- int, unsubscribed <-chan u.Never) {

			var count int

			drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   nil, // Not used
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						sendValue: func(chan<- T, <-chan u.Never, T) u.SelectResult {
							count++
							return u.Selection(u.SelectDone(unsubscribed), u.SelectSend(downstream, count))
						},
					}
				},
			})
		})
	}
}
