package rx

import "slices"

func Of[T any](values ...T) Observable[T] {
	return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan Never) {
		for value := range slices.Values(values) {
			if Selection(SelectDone(done), SelectSend(valuesOut, value)) {
				return
			}
		}
	})
}
