package rx

import u "alanpinder.com/rxgo/v2/utils"

func Filter[T any](filter func(T) bool) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:    source,
				valuesOut: valuesOut,
				errorsOut: errorsOut,
				done:      done,
				afterSelectionHandler: func(valueMsg *u.SelectReceiveMessage[T], errorMsg *u.SelectReceiveMessage[error]) AfterSelectionResult {

					if valueMsg.HasValue && !filter(valueMsg.Value) {
						return DropMessage
					}

					return ContinueMessage
				},
			})
		})
	}
}
