package rx

func Filter[T any](filter func(T) bool) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan Never) {

			valuesIn, errorsIn, unsubscribe := source.Subscribe()
			defer unsubscribe()

			for {

				var isValue, isError bool
				var value T
				var err error

				if Selection(SelectDone(done), SelectMustReceive(valuesIn, &isValue, &value), SelectMustReceive(errorsIn, &isError, &err)) {
					return
				}

				if isValue && !filter(value) {
					continue
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
