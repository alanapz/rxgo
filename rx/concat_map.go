package rx

func ConcatMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(valuesOut chan<- U, errorsOut chan<- error, done <-chan Never) {
			drainObservable(drainObservableArgs[T]{
				source:    source,
				valuesOut: nil, // Not used
				errorsOut: errorsOut,
				done:      done,
				valueHandler: func(_ chan<- T, done <-chan Never, value T) SelectResult {
					return drainObservable(drainObservableArgs[U]{source: projection(value), valuesOut: valuesOut, errorsOut: errorsOut, done: done}) == DoneResult
				},
			})
		})
	}
}
