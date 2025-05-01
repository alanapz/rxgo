package rx

func Take[T any](limit uint) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan Never) {

			var count uint

			drainObservable(drainObservableArgs[T]{
				source:    source,
				valuesOut: valuesOut,
				errorsOut: errorsOut,
				done:      done,
				afterSelectionHandler: func(valueMsg *SelectReceiveMessage[T], errorMsg *SelectReceiveMessage[error]) AfterSelectionResult {

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
