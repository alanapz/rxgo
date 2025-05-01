package rx

import (
	"slices"
)

func Concat[T any](sources ...Observable[T]) Observable[T] {
	return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan Never) {
		for source := range slices.Values(sources) {
			if drainObservable(drainObservableArgs[T]{source: source, valuesOut: valuesOut, errorsOut: errorsOut, done: done}) == DoneResult {
				return
			}
		}
	})
}
