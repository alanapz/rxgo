package rx

import u "alanpinder.com/rxgo/v2/utils"

func Take[T any](limit uint) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan u.Never) {

			var count uint

			drainObservable(drainObservableArgs[T]{
				source:    source,
				valuesOut: valuesOut,
				errorsOut: errorsOut,
				done:      done,
				afterSelectionHandler: func(valueMsg *u.SelectReceiveMessage[T], errorMsg *u.SelectReceiveMessage[error]) AfterSelectionResult {

					if valueMsg.HasValue {

						count++

						if count > limit {
							return ReturnContinue
						}
					}

					return ContinueMessage
				},
			})
		})
	}
}

func First[T any](source Observable[T]) Observable[T] {
	return Take[T](1)(source)
}
