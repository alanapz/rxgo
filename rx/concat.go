package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func Concat[T any](sources ...Observable[T]) Observable[T] {
	return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan u.Never) {
		for source := range slices.Values(sources) {
			if drainObservable(drainObservableArgs[T]{source: source, valuesOut: valuesOut, errorsOut: errorsOut, done: done}) == u.DoneResult {
				return
			}
		}
	})
}
