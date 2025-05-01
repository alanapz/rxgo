package rx

func Map[T any, U any](projection func(T) U) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(valuesOut chan<- U, errorsOut chan<- error, done <-chan Never) {

			valuesIn, errorsIn, unsubscribe := source.Subscribe()
			defer unsubscribe()

			for {

				var isValue, isError bool
				var value T
				var err error

				if Selection(SelectDone(done), SelectMustReceive(valuesIn, &isValue, &value), SelectMustReceive(errorsIn, &isError, &err)) {
					return
				}

				if isValue && Selection(SelectDone(done), SelectSend(valuesOut, projection(value))) {
					return
				}

				if isError && Selection(SelectDone(done), SelectSend(errorsOut, err)) {
					return
				}
			}
		})
	}
}

func ToAny[T any](input Observable[T]) Observable[any] {
	return Map(func(t T) any { return t })(input)
}

func FromAny[T any]() OperatorFunction[any, T] {
	return Map(func(t any) T { return t.(T) })
}
