package rx

import u "alanpinder.com/rxgo/v2/utils"

func Take[T any](limit uint) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(downstream chan<- T, unsubscribed <-chan u.Never) {

			var count uint

			drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   downstream,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						beforeSelection: func(*[]u.SelectItem) u.SelectResult {
							if count == limit {
								return u.DoneResult
							}
							return u.ContinueResult
						},
						onSelection: func(msg *u.SelectReceiveMessage[T]) AfterSelectionResult {

							if msg.HasValue {

								count++

								if count > limit {
									return StopAndContinueNext
								}
							}

							return ContinueMessage
						},
					}
				},
			})
		})
	}
}

func First[T any]() OperatorFunction[T, T] {
	return Take[T](1)
}
