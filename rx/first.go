package rx

func Take[T any](limit uint) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan Never) {

			valuesIn, errorsIn, unsubscribe := source.Subscribe()
			defer unsubscribe()

			var position uint

			for {

				var isValue, isError bool
				var value T
				var err error

				if Selection(SelectDone(done), SelectMustReceive(valuesIn, &isValue, &value), SelectMustReceive(errorsIn, &isError, &err)) {
					return
				}

				if isValue {

					position++

					if position >= limit {
						return
					}
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

func First[T any](source Observable[T]) Observable[T] {
	return Take[T](1)(source)
}
