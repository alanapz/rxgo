package rx

import u "alanpinder.com/rxgo/v2/utils"

func Map[T any, U any](projection func(T) U) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(valuesOut chan<- U, errorsOut chan<- error, done <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:    source,
				valuesOut: nil, // Not used
				errorsOut: errorsOut,
				done:      done,
				valueHandler: func(_ chan<- T, done <-chan u.Never, value T) u.SelectResult {
					return u.Selection(u.SelectDone(done), u.SelectSend(valuesOut, projection(value)))
				},
			})
		})
	}
}

func ToAny[T any](input Observable[T]) Observable[any] {
	return Map(func(t T) any { return t })(input)
}

func FromAny[T any]() OperatorFunction[any, T] {
	return Map(func(t any) T { return t.(T) })
}
