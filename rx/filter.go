package rx

func Filter[T any](filter func(T) bool) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan Never) {
			drainObservable(drainObservableArgs[T]{
				source:    source,
				valuesOut: valuesOut,
				errorsOut: errorsOut,
				done:      done,
				afterSelectionHandler: func(valueMsg *SelectReceiveMessage[T], errorMsg *SelectReceiveMessage[error]) AfterSelectionResult {

					if valueMsg.HasValue && !filter(valueMsg.Value) {
						return DropMessage
					}

					return ContinueMessage
				},
			})
		})
	}
}
