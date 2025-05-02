package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func TakeUntil[T any, X any](notifier <-chan X) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:    source,
				valuesOut: valuesOut,
				errorsOut: errorsOut,
				done:      done,
				additionalSelections: func() []u.SelectItem {
					return u.Of(u.SelectDone(notifier))
				},
			})
		})
	}
}
