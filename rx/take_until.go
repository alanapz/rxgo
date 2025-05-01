package rx

import "alanpinder.com/rxgo/v2/ux"

func TakeUntil[T any, X any](notifier <-chan X) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan Never) {
			drainObservable(drainObservableArgs[T]{
				source:    source,
				valuesOut: valuesOut,
				errorsOut: errorsOut,
				done:      done,
				additionalSelections: func() []SelectItem {
					return ux.Of(SelectDone(notifier))
				},
			})
		})
	}
}
