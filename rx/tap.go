package rx

import u "alanpinder.com/rxgo/v2/utils"

func Tap[T any](tapValue func(T), tapError func(error)) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, unsubscribed <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:       source,
				valuesOut:    valuesOut,
				errorsOut:    errorsOut,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onSelection: func(valueMsg *u.SelectReceiveMessage[T], errorMsg *u.SelectReceiveMessage[error]) AfterSelectionResult {

							if tapValue != nil && valueMsg.HasValue {
								tapValue(valueMsg.Value)
							}

							if tapError != nil && errorMsg.HasValue {
								tapError(errorMsg.Value)
							}

							return ContinueMessage
						},
					}
				},
			})
		})
	}
}
