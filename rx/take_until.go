package rx

func TakeUntil[T any, X any](notifier <-chan X) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan Never) {

			valuesIn, errorsIn, unsubscribe := source.Subscribe()

			defer unsubscribe()

			for {

				var isValue, isError bool
				var value T
				var err error

				if Selection(SelectDone(done), SelectDone(notifier), SelectMustReceive(valuesIn, &isValue, &value), SelectMustReceive(errorsIn, &isError, &err)) {
					return
				}

				if isValue && Selection(SelectDone(done), SelectSend(valuesOut, value)) {
					return
				}

				if isError && Selection(SelectDone(done), SelectSend(errorsOut, err)) {
					return
				}
			}
		})
	}
}
