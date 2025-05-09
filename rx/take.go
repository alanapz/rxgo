package rx

import u "alanpinder.com/rxgo/v2/utils"

func Take[T any](limit uint) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, unsubscribed <-chan u.Never) {

			var count uint

			drainObservable(drainObservableArgs[T]{
				source:       source,
				valuesOut:    valuesOut,
				errorsOut:    errorsOut,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onSelection: func(valueMsg *u.SelectReceiveMessage[T], errorMsg *u.SelectReceiveMessage[error]) AfterSelectionResult {

							if valueMsg.HasValue {

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
