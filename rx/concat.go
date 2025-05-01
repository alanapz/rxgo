package rx

import (
	"slices"
)

func Concat[T any](sources ...Observable[T]) Observable[T] {
	return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan Never) {

		for source := range slices.Values(sources) {

			if drainObservable(source, valuesOut, errorsOut, done) {
				return
			}
		}
	})
}
